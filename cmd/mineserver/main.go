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

type mineRoom struct {
	hub       *server.Hub
	mineCh    chan minesweeper.MineInputMsg
	setModeCh chan minesweeper.SetModeMsg
	done      chan struct{}
}

func startMineRoom(
	id string,
	stopCh <-chan struct{},
	saveResult func(minesweeper.MatchResult) error,
) *mineRoom {
	// 每个房间各自拥有独立的一套通道、Hub 和游戏循环，互不共享内存状态。
	broadcastCh := make(chan []byte, 64)
	joinCh := make(chan snake.JoinMsg, 4)
	leaveCh := make(chan snake.LeaveMsg, 4)
	startCh := make(chan snake.StartMsg, 4)
	gameOverCh := make(chan uint8, 1)
	mineCh := make(chan minesweeper.MineInputMsg, 32)
	setModeCh := make(chan minesweeper.SetModeMsg, 4)

	done := make(chan struct{})
	go func() {
		defer close(done)
		minesweeper.Loop(mineCh, joinCh, leaveCh, startCh, setModeCh,
			stopCh, broadcastCh, gameOverCh, saveResult)
	}()

	// 复用贪吃蛇的连接管理器，由它转发本房间的扫雷消息。
	inputCh := make(chan snake.InputMsg, 1)
	setTimeCh := make(chan snake.SetTimeMsg, 1)
	snakeModeCh := make(chan snake.SetModeMsg, 1)
	hub := server.NewHub(broadcastCh, inputCh, joinCh, leaveCh,
		startCh, setTimeCh, snakeModeCh)
	go hub.Run()

	go func() {
		for playerID := range gameOverCh {
			log.Printf("扫雷[房间%s]：对局结束，相关玩家编号=%d", id, playerID)
		}
	}()

	return &mineRoom{
		hub:       hub,
		mineCh:    mineCh,
		setModeCh: setModeCh,
		done:      done,
	}
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

	stopCh := make(chan struct{})

	// 只有配置了数据库才启用结算回调。每局写入最多等待 3 秒。
	var saveResult func(minesweeper.MatchResult) error
	if db != nil {
		saveResult = func(result minesweeper.MatchResult) error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return saveMatchResult(ctx, db, result)
		}
	}

	// 创建 rooms := map[string]*mineRoom{ ... }，分别用 startMineRoom("1", ...) 和 startMineRoom("2", ...) 启动两个房间
	rooms := map[string]*mineRoom{
		"1": startMineRoom("1", stopCh, saveResult),
		"2": startMineRoom("2", stopCh, saveResult),
	}
	// TODO(6)：把 "/ws" 路由改为调用 mineRoomWSHandler(rooms)
	http.HandleFunc("/ws", mineRoomWSHandler(rooms))
	// 数据库连接可用时，注册战绩列表和详情路由。
	if db != nil {
		http.HandleFunc("GET /matches/{gameID}", matchDetailHandler(db))
		http.HandleFunc("GET /matches", matchListHandler(db))
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
	//遍历 rooms，等待每个房间的 <-room.done 结束，再让 main 返回并关闭 db
	for _, room := range rooms {
		<-room.done
	}

	log.Println("扫雷服务已正常停止")

}

func mineRoomWSHandler(rooms map[string]*mineRoom) http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc, len(rooms))
	for id, room := range rooms {
		handlers[id] = mineWSHandler(room.hub, room.mineCh, room.setModeCh)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//  r.URL.Query().Get("room") 读取房间号；如果为空字符串 ""，默认设为 "1"
		roomnumbs := r.URL.Query().Get("room")
		if roomnumbs == "" {
			roomnumbs = "1"
		}
		// 从 handlers 中查找该房间对应的处理函数（handler, ok := handlers[...]）
		handler, ok := handlers[roomnumbs]
		// ：如果 !ok（未知房间号），调用 http.Error(w, "unknown room", http.StatusBadRequest) 并 return
		if !ok {
			http.Error(w, "unknown room", http.StatusBadRequest)
			return
		}
		// 调用查到的 handler(w, r) 把连接交给对应房间
		handler(w, r)
	}
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
