package minesweeper

import (
	"testing"
	"time"

	"github.com/siyad01/multiplayer-terminal-games/internal/protocol"
	"github.com/siyad01/multiplayer-terminal-games/internal/snake"
)

func TestRunOnceRejectsReusedConnection(t *testing.T) {
	lateAccepted, endedEarly, result := exerciseReusedConnection(t)
	t.Logf("lateAccepted=%v endedEarly=%v result=%+v", lateAccepted, endedEarly, result)

	// TODO(1)：检查延迟连接是否未被接受。
	if lateAccepted {
		t.Fatalf("不应该接受中途加入")
	}
	// TODO(2)：检查比赛是否在 Alice 自己离开之前结束。
	if endedEarly {
		t.Fatalf("对局在离开前结束了")
	}
	// TODO(3): 检查返回的快照是否存在、是否被取消、以及是否只有
	// 爱丽丝和鲍勃。爱丽丝的旗帜应该只反映她自己的有效旗帜。
	// 将下面的占位符替换为您的业务主张。
	if result == nil {
		t.Fatalf("结果为空了")
	}
	if !result.Cancel {
		t.Fatalf("比赛没有被正常取消")
	}
	if len(result.PlayerResult) != 2 {
		t.Fatalf("比赛玩家人数不对 期望是2,实际是%d", len(result.PlayerResult))
	}

	var foundAlice, foundBob bool
	for _, p := range result.PlayerResult {
		switch p.Name {
		case "Alice": //
			foundAlice = true
			if p.Flags != 1 {
				t.Fatalf("alice的插旗数不对,期望是1,实际是 %d", p.Flags)
			}
		case "Bob":
			foundBob = true
		default:
			t.Fatalf("有别的玩家错误加入: %+v", p)
		}
	}

	if !foundAlice || !foundBob {
		t.Fatalf("实际玩家不对: foundAlice=%v, foundBob=%v", foundAlice, foundBob)
	}
}

// AI-assisted fixture: drives the real runOnce through channels. Business
// assertions belong in TestRunOnceRejectsReusedConnection above.
func exerciseReusedConnection(t *testing.T) (bool, bool, *MatchResult) {
	t.Helper()
	const timeout = 3 * time.Second
	mineCh := make(chan MineInputMsg, 8)
	joinCh := make(chan snake.JoinMsg)
	// Unbuffered leaves make the two leave sends a processing boundary:
	// runOnce must finish Carol's branch before receiving Bob's leave.
	leaveCh := make(chan snake.LeaveMsg)
	startCh := make(chan snake.StartMsg, 1)
	setModeCh := make(chan SetModeMsg, 1)
	stopCh := make(chan struct{})
	broadcastCh := make(chan []byte, 64)
	gameOverCh := make(chan uint8, 1)
	done := make(chan *MatchResult, 1)
	exited := make(chan struct{})
	go func() {
		_, result := runOnce(mineCh, joinCh, leaveCh, startCh, setModeCh,
			stopCh, broadcastCh, gameOverCh, "fixed-roster")
		done <- result
		close(exited)
	}()
	t.Cleanup(func() {
		close(stopCh)
		select {
		case <-exited:
		case <-time.After(timeout):
			t.Error("fixture: runOnce did not stop")
		}
	})

	waitBroadcast := func(want protocol.MsgType) {
		t.Helper()
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		for {
			select {
			case data := <-broadcastCh:
				if len(data) > 0 && protocol.MsgType(data[0]) == want {
					return
				}
			case <-timer.C:
				t.Fatalf("fixture: timed out waiting for message %d", want)
			}
		}
	}
	join := func(id uint8, session uint64, name string) bool {
		t.Helper()
		reply := make(chan bool, 1)
		select {
		case joinCh <- snake.JoinMsg{PlayerID: id, SessionID: session, Name: name, ResultChan: reply}:
		case <-time.After(timeout):
			t.Fatal("fixture: join was not consumed")
		}
		select {
		case accepted := <-reply:
			return accepted
		case <-time.After(timeout):
			t.Fatal("fixture: no join reply")
			return false
		}
	}
	sendLeave := func(id uint8, session uint64) {
		t.Helper()
		select {
		case leaveCh <- snake.LeaveMsg{PlayerID: id, SessionID: session}:
		case <-time.After(timeout):
			t.Fatal("fixture: leave was not consumed")
		}
	}

	if !join(0, 1, "Alice") || !join(1, 2, "Bob") {
		t.Fatal("fixture: could not establish the initial two-player lobby")
	}
	waitBroadcast(protocol.MsgLobbyState)
	waitBroadcast(protocol.MsgLobbyState)
	setModeCh <- SetModeMsg{PlayerID: 0, SessionID: 1, Mode: ModeCooperative, Difficulty: "easy"}
	waitBroadcast(protocol.MsgLobbyState)
	startCh <- snake.StartMsg{PlayerID: 0, SessionID: 1}
	waitBroadcast(protocol.MsgMineState)

	// Carol represents a new connection reusing ID 0 while Alice's leave
	// notification has not yet been consumed. Hub timing is tested separately.
	lateAccepted := join(0, 3, "Carol")
	mineCh <- MineInputMsg{PlayerID: 0, SessionID: 3, Action: ActionFlag, X: 0, Y: 0}
	mineCh <- MineInputMsg{PlayerID: 1, SessionID: 2, Action: ActionCursor, X: 2, Y: 2}
	// FIFO on mineCh ensures Carol's input was processed before this cursor.
	waitBroadcast(protocol.MsgMineCursors)
	sendLeave(0, 3)
	sendLeave(1, 2)

	// Alice should still be able to flag. A premature result is captured for
	// the learner's assertion, rather than being treated as a fixture timeout.
	mineCh <- MineInputMsg{PlayerID: 0, SessionID: 1, Action: ActionFlag, X: 1, Y: 0}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case result := <-done:
			return lateAccepted, true, result
		case data := <-broadcastCh:
			if len(data) == 0 || protocol.MsgType(data[0]) != protocol.MsgMineState {
				continue
			}
			sendLeave(0, 1)
			select {
			case result := <-done:
				return lateAccepted, false, result
			case <-timer.C:
				t.Fatal("fixture: no result after final leave")
			}
		case <-timer.C:
			t.Fatal("fixture: neither Alice's update nor a result arrived")
		}
	}
}
