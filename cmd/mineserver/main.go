// cmd/mineserver/main.go
// 扫雷服务独立于贪吃蛇服务运行。
// 扫雷监听 :8081，避免与使用 :8080 的贪吃蛇服务冲突。

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/siyad01/multiplayer-terminal-games/internal/minesweeper"
	"github.com/siyad01/multiplayer-terminal-games/internal/protocol"
	"github.com/siyad01/multiplayer-terminal-games/internal/server"
	"github.com/siyad01/multiplayer-terminal-games/internal/snake"
)

func probeMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 探测失败时，关闭刚打开的数据库句柄并返回错误。
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func main() {
	var db *sql.DB
	dsn := os.Getenv("MINESWEEPER_MYSQL_DSN")
	if dsn != "" {
		var err error
		db, err = probeMySQL(dsn)
		if err != nil {
			log.Fatalf("扫雷服务：MySQL 连接探测失败：%v", err)
		}
		defer db.Close()
	}

	// 游戏循环和连接管理器共用的通道。
	broadcastCh := make(chan []byte, 64)
	joinCh := make(chan snake.JoinMsg, 4)
	leaveCh := make(chan snake.LeaveMsg, 4)
	startCh := make(chan snake.StartMsg, 4)
	gameOverCh := make(chan uint8, 1)
	stopCh := make(chan struct{})
	mineCh := make(chan minesweeper.MineInputMsg, 32)
	setModeCh := make(chan minesweeper.SetModeMsg, 4)

	// 只有配置了数据库才启用结算回调。每局写入最多等待 3 秒。
	var saveResult func(minesweeper.MatchResult) error
	if db != nil {
		saveResult = func(result minesweeper.MatchResult) error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return saveMatchResult(ctx, db, result)
		}
	}
	// 启动 Loop 时传入回调，并在 Loop 退出时发出完成信号。
	loopDone := make(chan struct{})
	go func() {
		defer close(loopDone)
		minesweeper.Loop(mineCh, joinCh, leaveCh, startCh, setModeCh,
			stopCh, broadcastCh, gameOverCh, saveResult)
	}()

	// 复用贪吃蛇的连接管理器，由它转发扫雷消息。
	inputCh := make(chan snake.InputMsg, 1) // 扫雷不用，但创建连接管理器时需要传入。
	setTimeCh := make(chan snake.SetTimeMsg, 1)
	snakeModeCh := make(chan snake.SetModeMsg, 1) // 扫雷不用，但创建连接管理器时需要传入。
	hub := server.NewHub(broadcastCh, inputCh, joinCh, leaveCh,
		startCh, setTimeCh, snakeModeCh)
	go hub.Run()

	go func() {
		for id := range gameOverCh {
			log.Printf("扫雷：对局结束，相关玩家编号=%d", id)
		}
	}()

	http.HandleFunc("/ws", mineWSHandler(hub, mineCh, setModeCh))
	// 数据库连接可用时，注册 GET /matches/{gameID} 战绩详情路由。
	if db != nil {
		http.HandleFunc("GET /matches/{gameID}", matchDetailHandler(db))
		// TODO(A3)：分页处理函数写好后，在这里注册 GET /matches。
	}

	srv := &http.Server{Addr: ":8081", Handler: nil}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("正在关闭扫雷服务……")
		close(stopCh)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("扫雷服务关闭失败：%v", err)
		}
	}()

	log.Printf("扫雷服务正在监听 :8081")
	log.Printf("启动客户端：go run cmd/client/main.go -name yourname")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	// 正常停机时等 Loop 结束，再让 main 返回并关闭 db。
	<-loopDone
	log.Println("扫雷服务已正常停止")

}

func mineWSHandler(
	hub *server.Hub,
	mineCh chan minesweeper.MineInputMsg,
	setModeCh chan minesweeper.SetModeMsg,
) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "anon"
		}
		if len(name) > protocol.MaxNameLen {
			name = name[:protocol.MaxNameLen]
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("升级为 WebSocket 连接失败：%v", err)
			return
		}

		c := hub.NewClient(name)
		hub.Register(c)

		// 把连接管理器发来的消息写给 WebSocket 客户端。
		go func() {
			defer conn.Close()
			for data := range c.SendChan() {
				if err := conn.WriteMessage(
					websocket.BinaryMessage, data); err != nil {
					return
				}
			}
		}()

		// 读取 WebSocket 客户端的消息，再转给扫雷游戏循环。
		go func() {
			defer func() {
				hub.Unregister(c)
				conn.Close()
			}()
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					return
				}

				msgType, err := protocol.ParseType(data)
				if err != nil {
					continue
				}

				switch msgType {
				case protocol.MsgMineReveal:
					if len(data) >= 3 {
						mineCh <- minesweeper.MineInputMsg{
							PlayerID:  c.ID(),
							SessionID: c.SessionID(),
							Action:    minesweeper.ActionReveal,
							X:         data[1],
							Y:         data[2],
						}
					}
				case protocol.MsgMineFlag:
					if len(data) >= 3 {
						mineCh <- minesweeper.MineInputMsg{
							PlayerID:  c.ID(),
							SessionID: c.SessionID(),
							Action:    minesweeper.ActionFlag,
							X:         data[1],
							Y:         data[2],
						}
					}
				case protocol.MsgMineCursor:
					if len(data) >= 4 {
						mineCh <- minesweeper.MineInputMsg{
							PlayerID:  c.ID(),
							SessionID: c.SessionID(),
							Action:    minesweeper.ActionCursor,
							X:         data[2],
							Y:         data[3],
						}
					}
				case protocol.MsgStartGame:
					hub.ForwardStart(c.ID(), c.SessionID())
				case protocol.MsgSetMineMode:
					mode, diff, err := protocol.ParseSetMineMode(data)
					if err != nil {
						continue
					}
					setModeCh <- minesweeper.SetModeMsg{
						PlayerID:   c.ID(),
						SessionID:  c.SessionID(),
						Mode:       minesweeper.Mode(mode),
						Difficulty: diff,
					}
				case protocol.MsgLeave:
					return
				case protocol.MsgPlayAgain:
					hub.Rejoin(c)
				}
			}
		}()
	}
}
