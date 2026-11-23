// internal/minesweeper/loop.go
// 扫雷循环由事件驱动，无需定时器轮询。
// 与贪吃蛇复用连接管理器和消息通道。

package minesweeper

import (
	cryptorand "crypto/rand"
	"log"
	"math/rand"
	"time"

	"github.com/siyad01/multiplayer-terminal-games/internal/protocol"
	"github.com/siyad01/multiplayer-terminal-games/internal/snake"
)

// MineInputMsg 携带一条玩家操作。
type MineInputMsg struct {
	PlayerID  uint8
	SessionID uint64
	Action    MineAction
	X, Y      uint8
}

type MineAction uint8

const (
	ActionReveal MineAction = iota
	ActionFlag
	ActionCursor
)

// Loop 运行扫雷游戏循环。
// 复用 snake.JoinMsg、snake.LeaveMsg 和 snake.StartMsg 通道。
func Loop(
	mineCh <-chan MineInputMsg,
	joinCh <-chan snake.JoinMsg,
	leaveCh <-chan snake.LeaveMsg,
	startCh <-chan snake.StartMsg,
	setModeCh <-chan SetModeMsg,
	stopCh <-chan struct{},
	broadcastCh chan<- []byte,
	gameOverCh chan<- uint8,
	saveResult func(MatchResult) error, // 可选的结算回调。
) {
	var results []MatchResult

	for {

		matchID := cryptorand.Text()
		done, result := runOnce(
			mineCh, joinCh, leaveCh, startCh, setModeCh,
			stopCh, broadcastCh, gameOverCh, matchID,
		)
		if result != nil {
			results = append(results, *result)
			log.Printf("扫雷：已保存对局结果到内存，对局编号=%s，累计=%d 局", result.GameID, len(results))
			// 有回调时同步保存此结果；失败记录对局编号和错误，然后继续下一局。
			if saveResult != nil {
				if err := saveResult(*result); err != nil {
					log.Printf("扫雷：对局 %s 写入数据库失败：%v", result.GameID, err)
				}
			}
		}
		if done {
			return
		}
		log.Println("扫雷：准备返回大厅")
	}
}

type SetModeMsg struct {
	PlayerID   uint8
	SessionID  uint64
	Mode       Mode
	Difficulty string
}

func runOnce(
	mineCh <-chan MineInputMsg,
	joinCh <-chan snake.JoinMsg,
	leaveCh <-chan snake.LeaveMsg,
	startCh <-chan snake.StartMsg,
	setModeCh <-chan SetModeMsg,
	stopCh <-chan struct{},
	broadcastCh chan<- []byte,
	gameOverCh chan<- uint8,
	matchID string,
) (bool, *MatchResult) {
	log.Println("扫雷：大厅已开放")

	lobbyPlayers := make(map[uint8]string)
	playerSessions := make(map[uint8]uint64)
	lobbyIDs := []uint8{}
	var hostID uint8
	hasHost := false
	gameMode := ModeSolo
	difficulty := "easy"

	broadcast := func(msg []byte) {
		select {
		case broadcastCh <- msg:
		default:
		}
	}

	addToLobby := func(id uint8, name string) bool {
		if len(lobbyPlayers) >= protocol.MaxPlayers {
			return false
		}
		if _, exists := lobbyPlayers[id]; exists {
			return true
		}
		lobbyPlayers[id] = name
		lobbyIDs = append(lobbyIDs, id)
		if !hasHost {
			hostID = id
			hasHost = true
		}
		return true
	}

	broadcastLobby := func() {
		ids := make([]uint8, 0, len(lobbyIDs))
		names := make([]string, 0, len(lobbyIDs))
		for _, id := range lobbyIDs {
			if n, ok := lobbyPlayers[id]; ok {
				ids = append(ids, id)
				names = append(names, n)
			}
		}
		hid := uint8(0xFF)
		if hasHost {
			hid = hostID
		}
		// 复用大厅状态消息，扫雷不设时长，所以时长填 0。
		broadcast(protocol.EncodeLobbyFull(hid, ids, names, 0,
			protocol.GameMode(gameMode+10))) // 模式编号加 10，以区别于贪吃蛇。
	}

	// ── 大厅 ─────────────────────────────────────────────────
	for {
		select {
		case msg := <-joinCh:
			ok := addToLobby(msg.PlayerID, msg.Name)
			msg.ResultChan <- ok
			if ok {
				playerSessions[msg.PlayerID] = msg.SessionID
				lobbyPlayers[msg.PlayerID] = msg.Name
				broadcastLobby()

			}

		case msg := <-leaveCh:
			// TODO：查询 playerSessions，取得保存的身份及是否存在
			// TODO：编号不存在，或保存身份与 msg.SessionID 不同，就 continue
			sessionID, ok := playerSessions[msg.PlayerID]
			if !ok || sessionID != msg.SessionID {
				continue
			}

			delete(lobbyPlayers, msg.PlayerID)
			delete(playerSessions, msg.PlayerID)
			for i, id := range lobbyIDs {
				if id == msg.PlayerID {
					lobbyIDs = append(lobbyIDs[:i], lobbyIDs[i+1:]...)
					break
				}
			}
			if hasHost && msg.PlayerID == hostID && len(lobbyIDs) > 0 {
				hostID = lobbyIDs[0]
			}
			if len(lobbyPlayers) == 0 {
				hasHost = false
			}
			broadcastLobby()

		case msg := <-setModeCh:
			sessionID, ok := playerSessions[msg.PlayerID]
			if !ok || sessionID != msg.SessionID {
				continue
			}

			if !hasHost || msg.PlayerID != hostID {
				continue
			}
			gameMode = msg.Mode
			difficulty = msg.Difficulty
			broadcastLobby()

		case msg := <-startCh:

			sessionID, ok := playerSessions[msg.PlayerID]
			if !ok || sessionID != msg.SessionID {
				continue
			}

			if !hasHost || msg.PlayerID != hostID {
				continue
			}
			if len(lobbyPlayers) == 0 {
				continue
			}
			log.Printf("扫雷：开局，模式编号=%d，难度=%s，玩家数=%d",
				gameMode, difficulty, len(lobbyPlayers))
			goto startGame

		case <-stopCh:
			return true, nil
		}
	}

startGame:
	diff := Difficulties[difficulty]
	if diff.Width == 0 {
		diff = Difficulties["easy"]
	}

	game := NewGame(matchID, diff, gameMode, rand.New(rand.NewSource(time.Now().UnixNano())))
	for _, id := range lobbyIDs {
		game.AddPlayer(id, lobbyPlayers[id])
		broadcast(protocol.EncodePlayerJoined(id, lobbyPlayers[id]))
	}

	// 发送初始棋盘状态。
	broadcast(protocol.EncodeMineState(
		game.Width, game.Height,
		game.VisibleCells(),
		game.Scores(),
		uint16(game.FlagsLeft),
	))

	log.Printf("扫雷：棋盘 %d×%d，地雷数=%d",
		game.Width, game.Height, game.MineCount)

	// ── 对局 ─────────────────────────────────────────────────
	// 只在收到事件时更新游戏。
	for {
		select {
		case msg := <-mineCh:
			sessionID, ok := playerSessions[msg.PlayerID]
			if !ok || sessionID != msg.SessionID {
				continue
			}

			if game.Over {
				continue
			}

			var result ActionResult
			var changes []protocol.MineChange

			switch msg.Action {
			case ActionReveal:
				result = game.Reveal(msg.PlayerID, msg.X, msg.Y)
				changes = result.Changes

			case ActionFlag:
				changes = game.Flag(msg.PlayerID, msg.X, msg.Y)

			case ActionCursor:
				game.MoveCursor(msg.PlayerID, msg.X, msg.Y)
				// 广播光标位置，每名玩家占 4 字节。
				broadcast(protocol.EncodeCursors(game.AllCursors()))
				continue
			}

			// 只广播发生变化的格子。
			if len(changes) > 0 {
				broadcast(protocol.EncodeMineUpdate(changes))
				// 同时广播分数。
				broadcast(protocol.EncodeMineState(
					game.Width, game.Height,
					game.VisibleCells(),
					game.Scores(),
					uint16(game.FlagsLeft),
				))
			}

			if result.Exploded {
				res := game.SnapshotResult()
				log.Printf("扫雷：玩家 %d 在 (%d,%d) 踩到地雷",
					msg.PlayerID, msg.X, msg.Y)
				broadcast(protocol.EncodeMineExplode(msg.PlayerID, msg.X, msg.Y))
				select {
				case gameOverCh <- msg.PlayerID:
				default:
				}
				time.Sleep(3 * time.Second)
				return false, &res
			}

			if result.Won {
				res := game.SnapshotResult()
				log.Println("扫雷：已翻开所有安全格，通关！")
				// 选出分数最高的玩家。
				var winner uint8 = 0xFF
				var best uint16 = 0
				for _, p := range game.Players {
					if p.Reveals > best {
						best = p.Reveals
						winner = p.ID
					}
				}
				broadcast(protocol.EncodeMineWin(game.Scores()))
				select {
				case gameOverCh <- winner:
				default:
				}
				time.Sleep(5 * time.Second)
				return false, &res
			}

		case msg := <-joinCh:
			// 对局进行中，拒绝新玩家加入本局。
			msg.ResultChan <- false

		case msg := <-leaveCh:
			sessionID, ok := playerSessions[msg.PlayerID]
			if !ok || sessionID != msg.SessionID {
				continue
			}

			game.RemovePlayer(msg.PlayerID)
			log.Printf("扫雷：玩家 %d 已离开本局", msg.PlayerID)

			// TODO：判断是否所有本局玩家都已离开
			if game.AllPlayersLeft() {

				// TODO：成立时调用取消方法，并记录包含 MatchId 的取消日志
				game.Cancel()
				res := game.SnapshotResult()
				gameID := game.MatchId
				log.Printf("扫雷：全员离开，对局已取消，对局编号=%v", gameID)
				// TODO：退出本次 runOnce，让外层 Loop 重新进入大厅
				return false, &res
			}
		case <-stopCh:
			return true, nil
		}
	}
}
