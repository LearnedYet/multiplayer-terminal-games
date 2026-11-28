package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/siyad01/multiplayer-terminal-games/internal/minesweeper"
	"github.com/siyad01/multiplayer-terminal-games/internal/protocol"
	"github.com/siyad01/multiplayer-terminal-games/internal/server"
	"github.com/siyad01/multiplayer-terminal-games/internal/snake"
)

const mineWSTestTimeout = 3 * time.Second

func TestMineCursorUsesConnectionIdentity(t *testing.T) {
	// AI-assisted network setup. The learner adds the identity assertion below.
	broadcastCh := make(chan []byte, 4)
	inputCh := make(chan snake.InputMsg, 1)
	joinCh := make(chan snake.JoinMsg, 4)
	leaveCh := make(chan snake.LeaveMsg, 4)
	startCh := make(chan snake.StartMsg, 1)
	setTimeCh := make(chan snake.SetTimeMsg, 1)
	snakeModeCh := make(chan snake.SetModeMsg, 1)
	mineCh := make(chan minesweeper.MineInputMsg, 4)
	setModeCh := make(chan minesweeper.SetModeMsg, 1)
	hub := server.NewHub(broadcastCh, inputCh, joinCh, leaveCh,
		startCh, setTimeCh, snakeModeCh)
	go hub.Run()

	testServer := httptest.NewServer(mineWSHandler(hub, mineCh, setModeCh))
	t.Cleanup(testServer.Close)
	var connections []*websocket.Conn
	t.Cleanup(func() {
		for _, conn := range connections {
			if err := conn.Close(); err != nil {
				t.Errorf("close test connection: %v", err)
			}
		}
		timer := time.NewTimer(mineWSTestTimeout)
		defer timer.Stop()
		for range connections {
			select {
			case <-leaveCh:
			case <-timer.C:
				t.Error("timed out waiting for test clients to unregister")
				return
			}
		}
	})

	connectPlayer := func(name string) (*websocket.Conn, uint8) {
		t.Helper()
		dialer := websocket.Dialer{HandshakeTimeout: mineWSTestTimeout}
		url := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws?name=" + name
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("connect %s: %v", name, err)
		}
		connections = append(connections, conn)
		if err := conn.SetReadDeadline(time.Now().Add(mineWSTestTimeout)); err != nil {
			t.Fatalf("set %s read deadline: %v", name, err)
		}
		messageType, welcome, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read %s welcome: %v", name, err)
		}
		if messageType != websocket.BinaryMessage || len(welcome) != 2 ||
			protocol.MsgType(welcome[0]) != protocol.MsgWelcome {
			t.Fatalf("unexpected %s welcome: type=%d data=%v", name, messageType, welcome)
		}
		playerID := welcome[1]
		select {
		case join := <-joinCh:
			join.ResultChan <- true
			if join.PlayerID != playerID || join.Name != name {
				t.Fatalf("unexpected join: got id=%d name=%s, want id=%d name=%s",
					join.PlayerID, join.Name, playerID, name)
			}
		case <-time.After(mineWSTestTimeout):
			t.Fatalf("timed out waiting for %s join", name)
		}
		return conn, playerID
	}

	aliceConn, aliceID := connectPlayer("alice")
	_, bobID := connectPlayer("bob")
	if aliceID == bobID {
		t.Fatalf("test needs two different player IDs, both were %d", aliceID)
	}
	if err := aliceConn.SetWriteDeadline(time.Now().Add(mineWSTestTimeout)); err != nil {
		t.Fatalf("set alice write deadline: %v", err)
	}
	if err := aliceConn.WriteMessage(websocket.BinaryMessage, protocol.EncodeCursor(bobID, 3, 4)); err != nil {
		t.Fatalf("send cursor claiming bob's ID: %v", err)
	}

	var msg minesweeper.MineInputMsg
	select {
	case msg = <-mineCh:
	case <-time.After(mineWSTestTimeout):
		t.Fatal("cursor message did not reach the game input channel")
	}
	if msg.Action != minesweeper.ActionCursor || msg.X != 3 || msg.Y != 4 {
		t.Fatalf("wrong cursor action or position: %+v", msg)
	}
	t.Logf("aliceID=%d bobID=%d received=%+v", aliceID, bobID, msg)

	if msg.PlayerID != aliceID {
		t.Fatalf("玩家编号不对 : 期望%d,实际%d", aliceID, msg.PlayerID)
	}
}

func TestMineRoomRoutingIsolatesRooms(t *testing.T) {
	// AI 辅助搭建双房间测试服务与连接辅助函数；核心断言由学习者在下面 TODO(1)-(3) 补全。
	stopCh := make(chan struct{})
	rooms := map[string]*mineRoom{
		"1": startMineRoom("1", stopCh, nil),
		"2": startMineRoom("2", stopCh, nil),
	}
	testServer := httptest.NewServer(mineRoomWSHandler(rooms))
	t.Cleanup(func() {
		testServer.Close()
		close(stopCh)
		for id, room := range rooms {
			select {
			case <-room.done:
			case <-time.After(mineWSTestTimeout):
				t.Errorf("房间 %s 循环未在停机后退出", id)
			}
		}
	})

	dialRoom := func(name, roomQuery string) (*websocket.Conn, <-chan []byte) {
		t.Helper()
		dialer := websocket.Dialer{HandshakeTimeout: mineWSTestTimeout}
		wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws?name=" + name
		if roomQuery != "" {
			wsURL += "&room=" + roomQuery
		}
		conn, _, err := dialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("连接 %s (room=%s) 失败：%v", name, roomQuery, err)
		}
		t.Cleanup(func() { _ = conn.Close() })
		ch := make(chan []byte, 32)
		go func() {
			defer close(ch)
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					return
				}
				ch <- data
			}
		}()
		return conn, ch
	}

	waitLobbyCount := func(ch <-chan []byte, wantCount uint8) {
		t.Helper()
		timer := time.NewTimer(mineWSTestTimeout)
		defer timer.Stop()
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					t.Fatalf("等待大厅人数 %d 时连接关闭", wantCount)
				}
				if len(data) >= 3 && protocol.MsgType(data[0]) == protocol.MsgLobbyState && data[2] == wantCount {
					return
				}
			case <-timer.C:
				t.Fatalf("超时未收到人数为 %d 的大厅状态", wantCount)
			}
		}
	}

	waitForMsgType := func(ch <-chan []byte, want protocol.MsgType) {
		t.Helper()
		timer := time.NewTimer(mineWSTestTimeout)
		defer timer.Stop()
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					t.Fatalf("等待消息类型 %#x 时连接关闭", byte(want))
				}
				if len(data) >= 1 && protocol.MsgType(data[0]) == want {
					return
				}
			case <-timer.C:
				t.Fatalf("超时未收到消息类型 %#x", byte(want))
			}
		}
	}

	ensureNoMsgType := func(ch <-chan []byte, forbidden protocol.MsgType) {
		t.Helper()
		timer := time.NewTimer(150 * time.Millisecond)
		defer timer.Stop()
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					t.Fatalf("检查无跨房串扰消息时连接意外关闭")
				}
				if len(data) >= 1 && protocol.MsgType(data[0]) == forbidden {
					t.Fatalf("发生跨房串扰：不应收到消息类型 %#x", byte(forbidden))
				}
			case <-timer.C:
				return
			}
		}
	}

	// 验证 1（未知房间返回 400；核心断言由用户写出期望状态码判断）。
	resp, err := http.Get(testServer.URL + "/ws?name=bad&room=99")
	if err != nil {
		t.Fatalf("请求未知房间失败：%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("错误房间号不应该被访问 房间号: %d", 99)
	}

	// 连接 1 号房两名玩家与 2 号房两名玩家，确认各自大厅人数都是 2（不串房）。
	r1a, r1aMsgs := dialRoom("r1a", "1")
	r1b, r1bMsgs := dialRoom("r1b", "1")
	waitLobbyCount(r1aMsgs, 2)
	waitLobbyCount(r1bMsgs, 2)

	r2a, r2aMsgs := dialRoom("r2a", "2")
	_, r2bMsgs := dialRoom("r2b", "2")
	waitLobbyCount(r2aMsgs, 2)
	waitLobbyCount(r2bMsgs, 2)

	// 1 号房房主切换合作模式并开局。
	if err := r1a.WriteMessage(websocket.BinaryMessage, protocol.EncodeSetMineMode(1, "easy")); err != nil {
		t.Fatalf("1 号房设置模式失败：%v", err)
	}
	if err := r1a.WriteMessage(websocket.BinaryMessage, protocol.EncodeStartGame()); err != nil {
		t.Fatalf("1 号房开局失败：%v", err)
	}

	// 验证 2：1 号房开局时，1 号房搭档收到棋盘而 2 号房玩家不会收到开局消息。
	waitForMsgType(r1bMsgs, protocol.MsgMineState)
	ensureNoMsgType(r2aMsgs, protocol.MsgMineState)

	// 2 号房也开局。
	if err := r2a.WriteMessage(websocket.BinaryMessage, protocol.EncodeSetMineMode(1, "easy")); err != nil {
		t.Fatalf("2 号房设置模式失败：%v", err)
	}
	if err := r2a.WriteMessage(websocket.BinaryMessage, protocol.EncodeStartGame()); err != nil {
		t.Fatalf("2 号房开局失败：%v", err)
	}
	waitForMsgType(r2aMsgs, protocol.MsgMineState)
	waitForMsgType(r2bMsgs, protocol.MsgMineState)

	// 1 号房全员断开取消本局并回到大厅。
	_ = r1a.Close()
	_ = r1b.Close()

	// 验证 3：1 号房全退重开期间，2 号房依然可以正常插旗并收到更新。
	if err := r2a.WriteMessage(websocket.BinaryMessage, protocol.EncodeFlag(0, 0)); err != nil {
		t.Fatalf("2 号房插旗失败：%v", err)
	}
	waitForMsgType(r2bMsgs, protocol.MsgMineUpdate)

	// 确认 1 号房清空后可重新接纳新玩家进入干净大厅。
	_, r1cMsgs := dialRoom("r1c", "1")
	waitLobbyCount(r1cMsgs, 1)
}
