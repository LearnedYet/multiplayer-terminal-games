<div align="center">

# 🎮 multiplayer-terminal-games（项目 A：多人合作扫雷服务端二次开发）

**基于 Go + WebSocket + MySQL 的多人合作扫雷服务端工程化改造**

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)

</div>

> 本仓库 Fork 自 [siyad01/multiplayer-terminal-games](https://github.com/siyad01/multiplayer-terminal-games)（Apache-2.0）。在保留原项目终端客户端、二进制协议与单 goroutine 状态所有权模型的基础上，围绕 **Minesweeper 合作模式（Cooperative / Easy）** 完成连接身份隔离、单局生命周期重构、MySQL 事务结算、HTTP 战绩查询、固定双房间隔离与本地浏览器试玩联调。

---

## 二次开发核心改造（A1 – A5）

- **输入信任边界与模式转发修复**
  - 修复 Snake 模式构造时遗漏保存 `setModeCh` 导致多人模式切换无效的缺陷，补充到达检查回归测试。
  - 修复扫雷操作信任客户端上报玩家编号的边界漏洞，入口统一强绑定服务端连接分配的 `c.ID()`，补充真 WebSocket 回归测试。
- **连接会话隔离与单对局生命周期**
  - 引入原子递增的连接级 `SessionID`，并在大厅与进行中对局校验 `(PlayerID, SessionID)`，阻断短 ID 回收复用引起的串局与迟到消息污染。
  - 实现开局锁定固定名单、玩家中途离开保留当前成绩并阻断后续操作、全员离开标记取消、结束生成不可变值拷贝快照（`SnapshotResult`）与随机唯一对局编号。
- **MySQL 事务结算与 HTTP 查询接口**
  - 新增 [`db/schema.sql`](db/schema.sql)（`match_results`、`match_players` 双表，InnoDB 联合主键与外键约束）。
  - 基于原生 `database/sql` 实现单事务原子落库与主键防重（1062）；同步写入设置 3 秒超时并记录错误日志，写入失败或超时后不中断后续对局循环。
  - 实现按对局查详情 `GET /matches/{gameID}` 与有界分页列表 `GET /matches?limit=...&offset=...`（限制每页 `1–50` 条、区分 400 / 404 / 500 状态），服务重启后已落库战绩可正常恢复查询。
- **固定双房间并发隔离（A5）**
  - 在服务端构建两个独立房间实例（`1` 号房、`2` 号房），各房维护独立的 `Hub`、消息通道与 `Loop` 协程，共用线程安全的 `*sql.DB` 落库与查询。
  - 在 `/ws` 接入层按 URL 查询参数 `?room=` 路由连接（未传默认进 1 号房，未知房号直接返回 `400`），停机时等待所有房间循环安全退出再释放数据库连接。
- **规则缺陷修复与回归保障**
  - 在指导下定位并修复合作模式跨玩家撤旗扣错玩家计数的归属缺陷（`Flag` 中先保存原插旗者 `FlaggedBy` 再清零扣减，配合 `uint16` 防下溢保护），补充 `TestCrossPlayerUnflagDecrementsOriginalOwner` 回归测试并在浏览器端实测验证。

### 当前范围与边界说明

- **已实现与验证**：固定双房间 Cooperative / Easy 隔离、对局生命周期、连接身份隔离、MySQL 事务结算与详情/分页查询、本机断库降级与重启回查验证；全库 `45` 项单元测试与集成测试通过。
- **暂未实现（不作夸大声明）**：动态多房间创建/销毁与跨房切换、Redis 缓存、用户注册登录体系及高并发压测指标。
- **分工与验证证据**：详细改动追踪、人工与 AI 协作边界及原始测试日志见 [`docs/CURRENT_STATE.md`](docs/CURRENT_STATE.md) 与 [`docs/EVIDENCE.md`](docs/EVIDENCE.md)。

---

## 本地验证与试玩说明

### 1. 运行全库测试（45 项）

```bash
go test -mod=readonly ./... -count=1 -v
```

### 2. 启动带 MySQL 持久化与浏览器试玩的服务

1. 在 MySQL 8.4 中创建数据库并执行 [`db/schema.sql`](db/schema.sql) 建表。
2. 在当前终端进程配置连接串环境变量 `MINESWEEPER_MYSQL_DSN`（例如 `user:password@tcp(127.0.0.1:13306)/go_a_learning`）。
3. 启动扫雷服务端，或使用本地浏览器试玩入口（详见 [`tools/playtest/README.md`](tools/playtest/README.md)）：

```powershell
# 启动扫雷服务端（监听 :8081）
go run ./cmd/mineserver

# 或启动带浏览器试玩页与战绩页面的本地联调入口（示例端口 8093）
pwsh -NoProfile -File .\tools\playtest\start.ps1 -Backend A3 -Port 8093
```

---

## 原上游说明（Upstream README）

<div align="center">

# 🎮 multiplayer-terminal-games

**Real-time multiplayer games that run entirely in your terminal**

*Snake and Minesweeper — playable over WebSockets from any terminal, anywhere on your network*

[![Build](https://github.com/siyad01/multiplayer-terminal-games/actions/workflows/build.yml/badge.svg)](https://github.com/siyad01/multiplayer-terminal-games/actions/workflows/build.yml)
[![Tests](https://img.shields.io/badge/Tests-29%20passing-22c55e)](#)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](go.mod)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows%20%7C%20WSL2-0ea5e9)](#install)

</div>

---

## Why this exists

Every multiplayer game framework assumes you have a browser, an app store, or a GPU.
This one assumes you have a terminal.

**multiplayer-terminal-games** is a WebSocket game server with two fully playable games —
Snake and Minesweeper — that render entirely in the terminal using raw mode.
No browser. No installation beyond a single binary. Works over SSH.

---

## Games

### Snake

Classic Snake with 10 levels, fruits, and poison.

| Mode | Description |
|------|-------------|
| Single Player | 10 levels, eat fruit to advance, avoid poison and walls |
| Last Standing | Multiplayer — last snake alive wins |
| Score Race | Multiplayer — most food eaten in the time limit wins |

- Levels 1-10: food count increases, poison added from level 2
- Different fruit emoji per level: 🍎 🍊 🍋 🍇 🍓 🍑 🍍 🥝 🍒 🫐
- Live scoreboard + timer (Score Race mode only)
- Host controls game duration with `+` / `-`

### Minesweeper

Classic Minesweeper with three difficulty levels and three play modes.

| Mode | Description |
|------|-------------|
| Solo | Classic single-player, clear the board without hitting a mine |
| Cooperative | Shared board — players work together to clear it |
| Competitive | Shared board — most cells revealed wins |

| Difficulty | Board | Mines |
|------------|-------|-------|
| Easy | 9×9 | 10 |
| Medium | 16×16 | 40 |
| Hard | 30×16 | 99 |

- Safe first click — first reveal never hits a mine
- Flood-fill reveal for empty cells
- Flag suspected mines with `F`
- Multiplayer cursors shown in real time

---

## Install

```bash
git clone https://github.com/siyad01/multiplayer-terminal-games
cd multiplayer-terminal-games
go build -o bin/snake-server    ./cmd/server
go build -o bin/mine-server     ./cmd/mineserver
go build -o bin/game-client     ./cmd/client
```

Or install directly:

```bash
go install github.com/siyad01/multiplayer-terminal-games/cmd/server@latest
go install github.com/siyad01/multiplayer-terminal-games/cmd/mineserver@latest
go install github.com/siyad01/multiplayer-terminal-games/cmd/client@latest
```

**Requirements:** Go 1.21+ — no other runtime, no external services.

---

## Quick start

```bash
# Terminal 1 — Snake server
./bin/snake-server

# Terminal 2 — Minesweeper server
./bin/mine-server

# Terminal 3+ — any number of players
./bin/game-client -name yourname

# Play on another machine on the same network
./bin/game-client -name yourname -server 192.168.x.x:8080
```

---

## Controls

### Snake

| Key | Action |
|-----|--------|
| `↑ ↓ ← →` or `WASD` | Move snake |
| `SPACE` / `ENTER` | Start game (host) |
| `+` / `-` | Adjust race timer (host, Score Race mode) |
| `M` | Cycle multiplayer mode (host) |
| `R` | Play again after game over |
| `Q` | Quit |

### Minesweeper

| Key | Action |
|-----|--------|
| `↑ ↓ ← →` or `WASD` | Move cursor |
| `ENTER` | Reveal cell |
| `F` | Flag / unflag mine |
| `TAB` | Cycle difficulty (menu) |
| `R` | Play again after game over |
| `Q` | Quit |

---
## Architecture

```text
multiplayer-terminal-games/
├── cmd/
│   ├── server/       # Snake WebSocket server (:8080)
│   ├── mineserver/   # Minesweeper server (:8081)
│   └── client/       # Unified terminal client
│
└── internal/
    ├── protocol/     # Binary wire format shared by client and server
    ├── snake/        # Snake game logic (pure, no networking)
    ├── minesweeper/  # Minesweeper logic, loop, and renderer
    ├── server/       # WebSocket hub and connection manager
    └── client/       # Terminal renderer and input handling
```

### Runtime Model

```text
[Client Terminal]
        │
        │ WebSocket (binary protocol)
        ▼
      [Hub]
        │
        ▼
[Game Loop Goroutine]
        ▲
        │
   ticker / events

Broadcasts game state updates every tick or player action.
```

### Key Design Decisions

* **Single-owner concurrency** — game state is owned by exactly one goroutine. All other goroutines communicate through typed channels. No mutexes are required for game state.
* **Binary protocol** — custom binary encoding rather than JSON. Every message begins with a 1-byte type header. A typical 30×20 board update is approximately 65 bytes.
* **Event-driven Minesweeper** — no ticker loop. The game wakes only when player actions occur.
* **Safe first click** — mines are generated after the first reveal, guaranteeing a safe 3×3 starting area.

---

## How the Protocol Works

Every message on the wire begins with a 1-byte message type.

### Client → Server

| Type   | Message           |
| ------ | ----------------- |
| `0x01` | Join              |
| `0x02` | Input (direction) |
| `0x03` | Leave             |
| `0x04` | StartGame         |
| `0x05` | SetTime           |
| `0x06` | PlayAgain         |
| `0x07` | NextLevel         |
| `0x08` | SetMode           |
| `0x09` | SetMineMode       |
| `0x20` | MineReveal        |
| `0x21` | MineFlag          |
| `0x22` | MineCursor        |

### Server → Client

| Type   | Message              |
| ------ | -------------------- |
| `0x10` | GameState            |
| `0x11` | PlayerJoined         |
| `0x12` | PlayerLeft           |
| `0x13` | GameOver             |
| `0x14` | Welcome              |
| `0x15` | LobbyState           |
| `0x16` | TimerUpdate          |
| `0x17` | LevelUp              |
| `0x18` | LevelComplete        |
| `0x19` | GameComplete         |
| `0x1A` | GameMode             |
| `0x30` | MineState            |
| `0x31` | MineUpdate (partial) |
| `0x32` | MineWin              |
| `0x33` | MineLose             |
| `0x34` | MineCursors          |

---

## Test Coverage

> Update this section to match the actual test count reported by `go test ./...`.

Example:

```text
36 tests across 4 packages — 0 failures

internal/protocol      12 tests  wire format round-trips
internal/snake         11 tests  game logic and food handling
internal/minesweeper    9 tests  mine placement, flood fill, flag toggling
internal/server         4 tests  hub behavior and player management
```

Run all tests:

```bash
go test ./...
```

---

## Comparison

| Feature | This project | typical game tutorial |
|---------|-------------|----------------------|
| Real multiplayer | ✅ WebSocket | ❌ single player only |
| Terminal rendering | ✅ tcell raw mode | ❌ browser only |
| Two games | ✅ Snake + Minesweeper | ❌ one game |
| Binary protocol | ✅ custom, documented | ❌ JSON |
| Tests | ✅ 29 passing | ❌ none |
| Single binary | ✅ no runtime | ❌ node_modules |

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache 2.0 — see [LICENSE](LICENSE)

---

<div align="center">

Built with Go · WebSockets · tcell · zero browser dependencies

**If you've ever wanted to play Snake over SSH, this is it. Give it a ⭐**

</div>