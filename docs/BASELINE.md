# 原始基线

## 1. 来源与环境

记录日期：2026-10-29，America/Los_Angeles（UTC-07:00）。本文件记录 Phase 0 实测结果，不把上游功能归为个人贡献。

| 项目 | 实测值 |
| --- | --- |
| 上游 | [siyad01/multiplayer-terminal-games](https://github.com/siyad01/multiplayer-terminal-games) |
| 工作目录 | `C:\Dev\Go\A`；最初为空目录，初查 Git 报 `fatal: not a git repository`，随后克隆到此处 |
| 原始提交 | `488dd313bc13de387d12a24f5d05a5c28bed5899`，`corrected README file`，提交时间 `2026-06-02T01:04:16Z` |
| Git 初始状态 | `main...origin/main`，无本地改动；`origin` fetch / push 均为上述仓库的 `.git` HTTPS 地址 |
| 系统 / Shell | Windows，系统 API 报告 `Microsoft Windows 10.0.26200`；PowerShell 7.6.5 Core |
| Git | `2.55.0.windows.2`；路径 `C:\Git\Git\cmd\git.exe` |
| Go | `go version go1.27.3 windows/amd64`；路径 `C:\Program Files\Go\bin\go.exe` |
| 模块 | `github.com/siyad01/multiplayer-terminal-games`，`go.mod` 要求 `go 1.26.3`；`GOMOD=C:\Dev\Go\A\go.mod`，无 `go.work` |
| Go 环境 | `GOOS=windows`，`GOARCH=amd64`，`GOTOOLCHAIN=auto`，`CGO_ENABLED=0` |
| 下载 / 校验 | `GOPROXY=https://proxy.golang.org,direct`，`GOSUMDB=sum.golang.org` |
| 缓存 | `GOPATH=C:\Users\c8618\go`，`GOCACHE=C:\Users\c8618\AppData\Local\go-build` |
| 许可证 | `LICENSE` 声明 Apache License 2.0，版权行为 `Copyright 2026 MultiPlayer Terminal Games Contributors`；保持原文件 |

已读取实际文件 `README.MD`、`CONTRIBUTING.MD`、`LICENSE`、`go.mod`、`.github/workflows/build.yml`，检查入口、核心模块、完整测试清单和代表性测试实现。上游没有 AGENTS.md。README 使用 `go build` / `go test`，没有额外服务、数据库或专用安装脚本。CI 配置有 Linux/macOS/Windows 矩阵，但本次未运行远程 CI。README 的 Go 1.21、测试数和 `-server` 示例不能替代源码和本次结果，差异见 [CURRENT_STATE.md](CURRENT_STATE.md)。

## 2. 依赖

使用现有 `go.mod` / `go.sum`，通过 `go mod download` 下载；未执行 `go get` / `go mod tidy`，未降低任何版本。以下是 go.mod 的 8 个 require 项；它们在上游均标注 indirect，但 tcell 和 websocket 被源码直接导入。

| 模块 | 锁定版本 | 用途 |
| --- | --- | --- |
| `github.com/gdamore/tcell/v2` | `v2.13.9` | 终端界面与键盘输入 |
| `github.com/gorilla/websocket` | `v1.5.3` | WebSocket 连接 |
| `github.com/gdamore/encoding` | `v1.0.1` | 终端字符编码支持 |
| `github.com/lucasb-eyer/go-colorful` | `v1.3.0` | 终端颜色支持 |
| `github.com/rivo/uniseg` | `v0.4.7` | Unicode 字符处理 |
| `golang.org/x/sys` | `v0.38.0` | 系统接口 |
| `golang.org/x/term` | `v0.37.0` | 终端支持 |
| `golang.org/x/text` | `v0.31.0` | 文本 / 编码支持 |

`go list -m all` 成功列出主模块及 17 个外部模块；除表内 8 项外，还有以下传递模块（不代表每个模块都用于当前 Windows 构建）：

```text
github.com/mattn/go-sixel v0.0.5
github.com/soniakeys/quant v1.0.0
github.com/yuin/goldmark v1.4.13
golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
golang.org/x/mod v0.29.0
golang.org/x/net v0.6.0
golang.org/x/sync v0.18.0
golang.org/x/tools v0.38.0
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
```

## 3. 实际执行的命令与结果

除初始目录探测外，工作目录均为 `C:\Dev\Go\A`。下表列出基线验证命令，文件读取及搜索不逐条展开。所有代理设置仅作用于对应 Git 命令或临时 PowerShell 子进程，未持久修改 Git / Go 配置，TLS 与模块校验未关闭。

| 命令 | 结果 |
| --- | --- |
| `git status --short --branch` / `git rev-parse HEAD` / `git log -1 --format="%H%n%ad%n%s" --date=iso-strict` / `git remote -v` | 获得上表提交、分支和来源；初始空目录中的 Git 探测失败，克隆后成功 |
| `git clone https://github.com/siyad01/multiplayer-terminal-games.git .` | 退出 1：`Recv failure: Connection was reset` |
| `Get-NetTCPConnection -State Listen -LocalPort 7897` | 确认本机代理端口当时在监听 |
| `git -c http.proxy=http://127.0.0.1:7897 -c http.sslVerify=true clone https://github.com/siyad01/multiplayer-terminal-games.git .` | 退出 0，成功克隆 |
| `git --version` / `go version` / `go env GOOS GOARCH GOMOD GOPATH GOPROXY GOTOOLCHAIN` / `go env GOMOD GOWORK GOSUMDB CGO_ENABLED CC GOCACHE` | 成功；环境见第 1 节 |
| `go mod download` | 退出 1：访问 proxy.golang.org 超时；关键错误见下方 |
| `$env:HTTPS_PROXY='http://127.0.0.1:7897'; $env:HTTP_PROXY='http://127.0.0.1:7897'; go mod download` | 退出 0，无输出，依赖下载成功 |
| `go mod verify` | 退出 0：`all modules verified` |
| `go list -m all` | 首次列举失败：上面 9 个传递模块的 `.info` 请求直连 proxy.golang.org 超时（`142.250.69.177:443: connectex`）；该次组合脚本未单独捕获 go 退出码，脚本最终退出 0 来自后续命令，不代表列举成功 |
| `$env:HTTPS_PROXY='http://127.0.0.1:7897'; $env:HTTP_PROXY='http://127.0.0.1:7897'; go list -m all` | 重试单独捕获 go 退出码为 0，完整模块版本见第 2 节 |
| `go test ./...` | 退出 0；4 个包测试成功，另外 4 个包无测试文件；原始输出见下方 |
| `go build ./...` | 退出 0，无输出 |
| `go build -o "$phase0Dir\snake-server.exe" ./cmd/server` | 退出 0，生成 Snake 服务端 |
| `go build -o "$phase0Dir\mine-server.exe" ./cmd/mineserver` | 退出 0，生成扫雷服务端 |
| `go build -o "$phase0Dir\game-client.exe" ./cmd/client` | 退出 0，生成统一客户端 |
| `git diff --exit-code HEAD -- go.mod go.sum LICENSE` / `git diff --exit-code HEAD` | 退出 0；所有上游受跟踪文件保持原样；未跟踪的新文档单独检查 |

构建时实际先执行：

```powershell
$phase0Dir = Join-Path $env:TEMP 'multiplayer-terminal-games-phase0-488dd313-20261029'
# 确认目录此前不存在后创建，用于仓库外临时产物。
New-Item -ItemType Directory -Path $phase0Dir | Out-Null
```

实际输出目录为 `C:\Users\c8618\AppData\Local\Temp\multiplayer-terminal-games-phase0-488dd313-20261029`。不需要 `go install ...@latest`，避免脱离固定提交构建。

首次依赖下载错误摘要（原始报错中的地址）：

```text
go: github.com/gdamore/encoding@v1.0.1:
Get "https://proxy.golang.org/github.com/gdamore/encoding/@v/v1.0.1.mod":
dial tcp [2607:f8b0:400a:803::2011]:443: connectex:
A connection attempt failed because the connected party did not properly respond after a period of time,
or established connection failed because connected host has failed to respond.
```

`go test ./...` 实际输出：

```text
?   github.com/siyad01/multiplayer-terminal-games/cmd/client              [no test files]
?   github.com/siyad01/multiplayer-terminal-games/cmd/mineserver          [no test files]
?   github.com/siyad01/multiplayer-terminal-games/cmd/server              [no test files]
?   github.com/siyad01/multiplayer-terminal-games/internal/client         [no test files]
ok  github.com/siyad01/multiplayer-terminal-games/internal/minesweeper    0.016s
ok  github.com/siyad01/multiplayer-terminal-games/internal/protocol       0.012s
ok  github.com/siyad01/multiplayer-terminal-games/internal/server         0.379s
ok  github.com/siyad01/multiplayer-terminal-games/internal/snake          0.293s
```

上述用时只是本次包测试输出，不是游戏服务端性能指标。静态统计：5 个 `_test.go` 文件，protocol 12、snake 11、minesweeper 9、server 4，共 36 个顶层测试；Snake 另有方向子测试。扫雷测试中存在条件性 `t.Skip`，本次普通输出未展开单个测试状态，因此不写“36 项逐一通过”或“零跳过”。现有测试覆盖协议往返、部分游戏规则、循环和 Hub 行为，没有终端客户端或完整 WebSocket 联机自动化测试。

## 4. 本次启动与联机验证

尝试用 `Start-Process -WindowStyle Hidden` 后台启动并重定向日志的组合命令，被工具策略在执行前拒绝，返回 `rejected: blocked by policy`，未给出更具体理由；没有据此认定程序启动失败。随后使用工具直接管理的独立终端会话启动成功。实际长运行命令如下，每条均在独立终端执行：

```powershell
& 'C:\Users\c8618\AppData\Local\Temp\multiplayer-terminal-games-phase0-488dd313-20261029\snake-server.exe'
& 'C:\Users\c8618\AppData\Local\Temp\multiplayer-terminal-games-phase0-488dd313-20261029\mine-server.exe'
& 'C:\Users\c8618\AppData\Local\Temp\multiplayer-terminal-games-phase0-488dd313-20261029\game-client.exe' -name phase0_alice
& 'C:\Users\c8618\AppData\Local\Temp\multiplayer-terminal-games-phase0-488dd313-20261029\game-client.exe' -name phase0_bob
```

客户端每轮退出后，以同样命令重新运行并选择另一款游戏。本机终端会话支持键盘输入及 tcell 输出，进行了实际双客户端冒烟。下面表内客户端路径用第 3 节的 `$phase0Dir` 简写。

| 检查 / 实际操作 | 结果与证据 |
| --- | --- |
| `Get-NetTCPConnection -State Listen -LocalPort 8080,8081` | 两端口监听；源码绑定 `:8080` / `:8081`，本次仅连接本机，未部署或修改防火墙 |
| `Invoke-WebRequest -Uri 'http://127.0.0.1:8080/' -TimeoutSec 5` | HTTP 200，正文 `multiplayer-terminal-games snake server` |
| `Invoke-WebRequest -Uri 'http://127.0.0.1:8081/ws' -TimeoutSec 5 -SkipHttpErrorCheck` | HTTP 400；普通 HTTP 未带 WebSocket Upgrade，属于预期拒绝，服务端记录 upgrade error |
| `& "$phase0Dir\game-client.exe" -h` | 退出 0，仅有 `-name string`，没有 `-server` |
| Snake：两个客户端均在菜单按下方向键一次、Enter，选择 Last Standing | 日志确认 alice=0、bob=1，均加入同一大厅；两端显示两名玩家，但标题为 SINGLE PLAYER |
| Snake：Alice 按 Enter 开局，Bob 按 W | 两端收到棋盘及后续更新；日志 `starting — mode=1 players=2`，随后 `ticking level 1 mode 1`。双人传输链路通过；Last Standing 模式选择失败（R1），不能把本次开局当作其规则验收通过 |
| Minesweeper：重启两客户端，菜单各按下方向键四次、Enter，选择 Cooperative / Easy；Alice 先加入 | 两端大厅均显示 Cooperative 与两名玩家 |
| Minesweeper：Alice 按 Enter 开局，Bob 按 Enter 翻格 | 日志 `starting mode=1 diff=easy players=2`、`board 9x9 mines=10`；两端均显示展开区域和 Bob 分数 44。这里只将该分数作为同步证据，随机棋盘不保证复跑相同分数 |
| 两轮客户端均按 Q；服务端分别按 Ctrl+C | 客户端均退出 0。服务端记录 `shutting down ...`、`stopped cleanly`，但工具终端会话退出码均为 1；只确认停机日志与端口释放，不宣称退出码为 0 |
| 停机后再次查询 8080/8081 | 剩余监听数为 0，本次启动进程已清理 |

**未验证**：Last Standing 正确胜负规则、Score Race、扫雷完整通关 / 所有模式难度、重连与异常断线、慢客户端、恶意输入、长时间运行、Race 检测、覆盖率、基准 / 压力测试、跨机器、Linux/macOS、远程 CI。模式缺陷是已复现失败，不属于“未验证”。

## 5. Windows 本地复验步骤（待用户执行，不代表已执行）

PowerShell 7，在仓库根目录准备本机程序；若默认网络失败，先确认本地代理，再只在当前 PowerShell 设置第 3 节的代理环境变量。

```powershell
Set-Location C:\Dev\Go\A
go version
go mod download
go test ./...
go build ./...
$localBin = Join-Path $env:TEMP 'multiplayer-terminal-games-local-check'
New-Item -ItemType Directory -Path $localBin -Force | Out-Null
go build -o "$localBin\snake-server.exe" ./cmd/server
go build -o "$localBin\mine-server.exe" ./cmd/mineserver
go build -o "$localBin\game-client.exe" ./cmd/client
```

打开四个 Windows Terminal / PowerShell 窗口，均先设置 `$localBin = Join-Path $env:TEMP 'multiplayer-terminal-games-local-check'`，分别执行：

```powershell
# 窗口 1
& "$localBin\snake-server.exe"
# 窗口 2
& "$localBin\mine-server.exe"
# 窗口 3
& "$localBin\game-client.exe" -name alice
# 窗口 4
& "$localBin\game-client.exe" -name bob
```

1. 先用 Snake Last Standing 重现 R1：Alice 先确认菜单，Bob 再加入；大厅应显示两名玩家，记录实际模式文字；Alice 按 Enter，观察日志与两端棋盘。当前基线预期会出现单人模式错误，不因能开局就判定多人模式通过。
2. 两客户端 Q 退出并重新运行，均选 Minesweeper Cooperative / Easy；Alice 先加入，等 Bob 出现在双方大厅后按 Enter 开局；Bob 按 Enter 翻格，检查 Alice 的棋盘与 Bob 分数是否同步；再按方向键 / F 检查光标与旗帜（本次未验证这两项）。
3. 完整胜负、重开或异常退出若需验收，另记每项输入、预期与实际结果；禁止用未执行步骤填“通过”。
4. 客户端 Q 退出，服务端 Ctrl+C；用 `Get-NetTCPConnection -State Listen -LocalPort 8080,8081 -ErrorAction SilentlyContinue` 检查监听释放，保留输出或截图。

两个客户端当前地址固定为 localhost，不使用 README 的 `-server` 参数。若端口已占用，先确认占用者；不要结束不属于本次验收的进程。

## 6. 关键模块与调用链（依据源码）

| 模块 | 职责与入口 |
| --- | --- |
| `cmd/server/main.go` | 创建 Snake 消息通道，启动 `snake.Loop`、`Hub.Run`，注册 `/ws`，监听 8080 |
| `cmd/mineserver/main.go` | 启动 `minesweeper.Loop` 并复用 Hub；在入口内处理扫雷 WebSocket 消息，监听 8081 |
| `cmd/client/main.go` | `main → showMainMenu → runSnake / runMine`；先选菜单再连接；接收网络消息、处理键盘并绘制 |
| `internal/server` | `ServeWS` 升级连接；`Client.readPump/writePump` 收发 Snake 消息；`Hub.Run` 管理客户端、分配 ID、广播 |
| `internal/protocol/messages.go` | 双方共用的二进制消息类型、Encode / Parse；首字节是类型，后面是对应载荷；最多 4 名玩家 |
| `internal/snake` | `loop.go` 用通道及 ticker 推进大厅和对局；`game.go` 处理移动、碰撞、分数、关卡及棋盘 |
| `internal/minesweeper` | `loop.go` 按事件推进；`game.go` 处理雷区与翻格 / 旗帜 / 光标；`render.go` 是客户端渲染器 |
| `internal/client` | Snake 的 tcell 渲染和键盘事件解析 |

Snake 输入链：`键盘 → client.ParseKey / EncodeInput → WebSocket → Client.readPump → ParseInput → Hub.ForwardInput → inputCh → snake.Loop → Game.SetDirection / Tick`。

Snake 状态返回链：`Game.Board / Scores → protocol.EncodeGameState → broadcastCh → Hub.Run → Client.send → writePump → WebSocket → 客户端 ParseGameState → Renderer.DrawGame`。模式链另走 `EncodeSetMode → readPump → ForwardSetMode → setModeCh → Loop`，目前在 Hub 字段初始化处中断（R1）。

扫雷链：`runMine 键盘 → EncodeReveal / Flag / Cursor → WebSocket → cmd/mineserver 的读 goroutine → mineCh → minesweeper.Loop → Game.Reveal / Flag / MoveCursor → 编码状态 / 增量 → broadcastCh → Hub → 入口内写 goroutine → 客户端解析 → minesweeper.Renderer`。

服务端游戏状态主要由各自的循环 goroutine 维护，其他 goroutine 通过 channel 传事件。channel 可理解为不同执行流程间传递消息的队列；Hub 管连接和转发，Game 管规则。扫雷仍复用 Snake 的加入 / 离开消息类型，当前并非完全独立的通用游戏框架。本阶段保留现有结构，避免在验证之前改变基线。

学习重点：Go module 与版本锁定、`cmd` / `internal` 分工、goroutine 与 channel、WebSocket 生命周期、二进制编解码，以及“构建成功 / 包测试成功 / 联机冒烟 / 规则正确”之间的证据差别。

## 7. 交付检查

新增文件只有 `AGENTS.md`、`docs/CURRENT_STATE.md`、本文件及 `docs/EVIDENCE.md`。采用“一份基线记录实测、一份状态记录风险、一份证据模板记录贡献”的分工，避免建立重复的架构文档。已核对文件范围、文档链接、阶段边界及通过 / 失败 / 未验证表述；上游受跟踪文件没有差异，HEAD 未变，临时进程已停止。网络问题已用进程级代理解决；R1 阻碍 Snake 多人模式的正确验收，不阻碍本阶段记录原始基线。Phase 0 交付后停止。
