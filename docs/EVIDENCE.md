# 贡献与测试证据

## 记录原则

- 上游来源、原始提交与原有能力见 [BASELINE.md](BASELINE.md)。上游实现及原有测试通过，不计作个人功能或性能成果。
- 每项贡献区分：用户独立完成的部分、上游复用部分、工具 / AI 辅助部分；保留可定位的提交、差异和可复现命令。
- 原始 Phase 0 为用户确定目标与约束、Codex 辅助执行环境验证并整理四份文档；没有个人业务代码实现，没有新增 Git 提交。原始运行结果保留在 BASELINE，本人后续复验另记如下，不重写历史基线。
- 下面的空白模板用于未来实际发生的工作。空白代表尚未填写，不代表 0、通过或已完成。
- 各批次保留执行当时的状态；A1.1 最新结论见本文的“学习收尾与交付”和个人贡献索引，早期“待验收”不代表当前仍未完成。

## 本人最小复验记录

- 批次：`manual-check-20261030-20261102`；执行对话：本项目 A 的 Phase 0 聊天。2026-10-30 开始，2026-11-02 根据补充截图收尾。
- 分工：用户亲手操作 Windows PowerShell / 终端；AI 提供分步指引、核对用户截图并记录。游戏与同步逻辑均为上游实现，本次没有编写核心代码，也不计为个人功能成果。
- 本轮只读核对 HEAD 仍为 `488dd313bc13de387d12a24f5d05a5c28bed5899`。未重新执行包测试、构建命令、端口检查或性能测试；用户的 `go run` 本身会编译并启动。

用户终端命令（每个长运行命令分别在独立窗口执行，各窗口先切换目录）：

```powershell
Set-Location C:\Dev\Go\A
go run ./cmd/mineserver
go run ./cmd/client -name alice
go run ./cmd/client -name bob
```

| 项目 | 实际证据与结论 |
| --- | --- |
| 同局 | 用户在指导下选择 Minesweeper Cooperative / Easy，并回报两端大厅均列出 alice、bob。 |
| 翻格同步 | [双端截图](evidence/manual-check-20261030-20261102/01-board-sync.png)中展开区域一致，两端记分牌均为 alice 49、bob 0；同步已复现。未单独核实具体触发翻格的玩家，不记录为 Bob 得分。 |
| 客户端退出 | [退出截图](evidence/manual-check-20261030-20261102/02-clients-returned.png)确认两端返回 PowerShell 提示符；图中有 `^C`，无法据此确认均通过 Q 退出，客户端退出码未采集。 |
| 途中失败与纠正 | 用户在原客户端窗口重复运行 `go run ./cmd/mineserver`，出现 [8081 绑定失败](evidence/manual-check-20261030-20261102/03-duplicate-server-bind-failure.png)：`Only one usage of each socket address ... is normally permitted`、`exit status 1`。指引已纠正：应切回正在运行的服务端窗口按 Ctrl+C。 |
| 服务端停机 | [11 月 2 日截图](evidence/manual-check-20261030-20261102/04-server-clean-shutdown.png)显示 `go run ./cmd/mineserver` 启动后，出现 `shutting down minesweeper server...`、`minesweeper server stopped cleanly` 并返回提示符。截图日志为 `2026/11/02 22:08:36` 启动、`22:08:48` 停机；这些时间按终端原样记录，未另核对终端时区。未采集退出码。 |

截图按用户原始附件复制到本批次独立目录，逐份核对 SHA256 一致，未覆盖原有证据。同步、客户端退出与后续服务端停机是分段观察，不宣称属于同一次连续联机；也未证明之前占用 8081 的具体进程如何退出。本人复验到此结束，未追加重测。

最小调用链：`客户端翻格输入 → WebSocket → 扫雷服务端游戏循环 → 广播状态 → 两端绘制`，具体函数见 BASELINE。学习结果为“指导下完成运行与观察”，尚不等于完全独立操作、掌握内部实现或能独立维护；游戏分数不是性能数据。本次仅更新 CURRENT_STATE 与本文档、留存截图，不修改业务代码或开启下一阶段。

## A1.1 编辑器准备（2026-11-03，已完成）

- 批次：`a11-editor-20261103`，执行对话为本项目 A 的 Phase 0 / A1.1 聊天；用户手动安装 VS Code Go 扩展，AI 只读排查并给出工具安装指引，未执行 gopls 安装或编写核心代码。
- [VS Code 原始失败日志](evidence/a11-editor-20261103/gopls-auto-install-failure.log)：扩展执行 `go install -v golang.org/x/tools/gopls@latest`，退出码 1，访问 `https://proxy.golang.org/golang.org/x/tools/gopls/@v/list` 直连超时。副本来自用户 Code 日志目录 `20261103T001058/window1/exthost/golang.go/Go.log`，保留本次独立批次，不覆盖既有证据。
- 本轮读取 `go version` 为 `go1.27.3 windows/amd64`，`GOBIN=C:\Users\c8618\go\bin`；预期路径下未找到 `gopls.exe`，PATH 也未找到该命令。
- 沙箱内网络及 CIM 检查受权限限制，未据此判定本机代理离线。后续获准只读检查：7897 处于监听状态；`Invoke-WebRequest -Uri 'https://proxy.golang.org/golang.org/x/tools/gopls/@latest' -Proxy 'http://127.0.0.1:7897' -TimeoutSec 10` 返回 HTTP 200，响应版本为 v0.24.0。此条是实测摘要，不是原始日志副本，也不证明 gopls 已安装。
- 用户按临时 HTTP_PROXY / HTTPS_PROXY 指引重试，截图显示已下载多项 gopls 依赖，但停在 `github.com/segmentio/encoding v0.5.4`，见[下载停滞截图](evidence/a11-editor-20261103/gopls-install-download-stall.png)（从附件复制，SHA256 一致）。只读检查见一个 `go.exe install golang.org/x/tools/gopls@latest` 进程（PID 17000），依赖下载临时 zip 文件持续为 0 字节，gopls.exe 仍不存在；未终止该进程。
- 后续网络检查摘要：`Invoke-WebRequest -Method Head -Uri 'https://proxy.golang.org/github.com/segmentio/encoding/@v/v0.5.4.zip' -Proxy 'http://127.0.0.1:7897' -TimeoutSec 8` 超时；对同一模块版本改用 `https://goproxy.cn/github.com/segmentio/encoding/@v/v0.5.4.zip` 并带 `-NoProxy`，返回 HTTP 200，Content-Length=13831857。这只验证了响应头，不代表完整 zip 已下载或安装成功。
- 用户随后提供[安装命令返回截图](evidence/a11-editor-20261103/gopls-install-returned.png)，可见此前代理重试及临时 `GOPROXY=https://goproxy.cn,direct`、`NO_PROXY=goproxy.cn` 的 `go install -v golang.org/x/tools/gopls@latest` 均已返回提示符；未逐条采集安装退出码，不单独归因哪次完成编译。截图副本 SHA256 一致。
- 安装后的实际核验 `& 'C:\Users\c8618\go\bin\gopls.exe' version`：输出 `golang.org/x/tools/gopls v0.24.0`，退出码 0；文件存在，大小 44775424 字节，安装已确认。另检查 `git diff --exit-code HEAD -- go.mod go.sum LICENSE cmd internal` 返回 0，项目源码及依赖未变。未重跑业务测试、未更改系统永久配置。
- 用户在重载与跳转检查指引后回复“能跳了”，据此记录 F12 功能由用户确认可用；AI 未直接操作或观察该次跳转。编辑器准备到此结束，回到 A1.1 的第一个测试编辑步骤；不将工具准备计为用户业务代码贡献。

## A1.1 本人修复前 Snake 观察（2026-11-03）

- 批次：`a11-snake-before-fix-20261103`；执行对话：本项目 A 当前 Phase 0 / A1.1 聊天。用户亲手操作 Windows PowerShell，AI 给出运行与观察指引、核对附件并整理文档；没有业务代码贡献。
- 目的与输入：在两个客户端选择 `Snake - Last Standing (multiplayer)`，由房主开局，观察服务端实际采用的模式。客户端按键过程未单独录像或截图；以下结果以用户按指引回传的截图为依据。
- 环境与指引命令：工作目录 `C:\Dev\Go\A`；窗口 1 执行 `go run ./cmd/server`，窗口 2 / 3 分别执行 `go run ./cmd/client -name alice` / `go run ./cmd/client -name bob`；每个窗口先执行 `Set-Location C:\Dev\Go\A`。两客户端菜单初始位置按 ↓ 一次、Enter 选择 Last Standing，进入大厅后由 alice 按 Enter 开局。服务端命令可在附件中直接看到，客户端启动命令未单独截图。

| 观察点 | 证据与结果 |
| --- | --- |
| 启动 | [启动截图](evidence/a11-snake-before-fix-20261103/01-server-started.png)显示 `snake server listening on :8080`、`game loop: lobby open`、`hub started`。 |
| 大厅 | [双端大厅截图](evidence/a11-snake-before-fix-20261103/02-two-player-lobby.png)中两端均列出 `[H] alice`、`bob`，标题均为 `SINGLE PLAYER - SNAKE`；alice 有开局提示，bob 等待房主。 |
| 开局 | [服务端日志截图](evidence/a11-snake-before-fix-20261103/03-started-mode-1.png)记录 alice 为玩家 0 / 房主、bob 为玩家 1；`2026/11/03 01:19:34 game loop: starting — mode=1 players=2`，随后 `ticking level 1 mode 1`。Last Standing 的预期值为 2，实际为 1（单人），本轮观察到模式不一致。时间按截图原样记录，未另核对终端时区。 |
| 退出 | [离开与停机截图](evidence/a11-snake-before-fix-20261103/04-clients-left-server-stopped.png)记录 01:24:09 客户端 1、01:24:11 客户端 0 发出 leave message，Hub 注销至 `total=0`，两连接 readPump / writePump 退出；01:24:33 出现 `shutting down snake server...`、`snake server stopped cleanly` 并返回 `PS C:\Dev\Go\A>`。用户回复“出现了”；各客户端提示符及具体退出按键未单独截图，未采集进程退出码。 |

四张用户原始附件已复制到本批次独立目录，逐份核对 SHA256 一致，没有覆盖旧证据。源码只读核对 `git diff --exit-code HEAD -- internal/server/hub.go internal/server/hub_test.go` 返回 0；未新增或执行回归测试，未修改业务代码，也未运行全库测试、构建检查或性能测试。本次手工现象不能替代后续同一回归测试的修复前失败与修复后通过证据；未验证完整游戏规则或用户独立调试能力。

## A1.1 参数与初始化相近练习（2026-11-04，说明待验收）

- 批次：`a11-initialization-exercise-20261104`；执行对话：本项目 A 当前 Phase 0 / A1.1 聊天。目标是区分传入参数、结构体字段与后续读取；本次是学习性改名，不计作新功能或 R1 修复成果。
- 示范与分工：AI 用上游正常 `startCh` 交接及 C# 成员赋值作示范，给出改名目标和定位，未提供练习的完整修改补丁；用户自行修改实际文件并回传截图，AI 只读核对、编译和整理证据。
- 实际改动：`internal/server/hub.go` 的 `NewHub` 参数 `setTimeCh` 改名为 `incomingTimeCh`，字段初始化改为 `setTimeCh: incomingTimeCh`；字段仍为 `Hub.setTimeCh`，后续 `ForwardSetTime` 仍读取 `h.setTimeCh`。同文件另有格式整理，未见其他行为修改。见[实际差异](evidence/a11-initialization-exercise-20261104/user-exercise.diff)与[用户截图](evidence/a11-initialization-exercise-20261104/user-initialization-edit.png)，截图副本 SHA256 与原附件一致。
- 交接链：`main 中的 setTimeCh → NewHub 参数 incomingTimeCh → Hub 字段 setTimeCh → ForwardSetTime 读取 h.setTimeCh`。这是已有初始化关系的改名练习，上游架构与方法逻辑不计为个人实现。
- 实际验证：Windows amd64、Go 1.27.3，工作目录 `C:\Dev\Go\A`，基准 HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；`git diff --check -- internal/server/hub.go internal/server/hub_test.go` 返回 0，`go build -mod=readonly ./internal/server` 无输出、退出码 0。命令、环境和退出码保存在[编译记录](evidence/a11-initialization-exercise-20261104/build.log)。依赖及许可证差异检查为空。
- 掌握状态：给定目标后的相近修改已通过实际文件检查及编译；用户本人对三个环节的说明尚待完成，第一学习块尚未验收。未运行测试、未验证设置时间的实际收发行为；现有模式测试草稿仍只有到达检查，R1 未修复。本轮未进入 channel 学习或追加测试断言。

## A1.1 参数与字段练习复核（2026-11-07）

- 批次：`a11-initialization-exercise-20261107-01`；执行对话：本项目 A 当前 Phase 0 / A1.1 聊天。用户编写练习代码，AI 提供整数示例、改名目标与定位提示，读取实际文件、编译及维护文档；未代写字段赋值。本次不计作新功能或 R1 修复。
- 用户操作与说明：在 `internal/server/hub.go` 的 `NewHub` 中将参数 `startCh` 改为 `incomingStartCh`，自行将初始化改为 `startCh: incomingStartCh`。用户回复“好了 不用改名，参数保存在左边的字段了”，能说明后续方法读取原字段，所以不随参数改名。准确表述为参数携带的值保存到字段；AI 已核对 `Hub.startCh` 及 `ForwardStart` 的 `h.startCh` 保持一致。
- 关键交接：`调用方 startCh → NewHub 参数 incomingStartCh → Hub 字段 startCh → ForwardStart 读取 h.startCh`。这是现有交接关系的学习练习，上游方法和架构仍属上游实现。此前 `incomingTimeCh` 练习保持原样。
- 实际验证：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`，HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；`git diff --check -- internal/server/hub.go` 与 `go build -mod=readonly ./internal/server` 均无输出、退出码 0，见[编译原始记录](evidence/a11-initialization-exercise-20261107-01/build.log)。[差异快照](evidence/a11-initialization-exercise-20261107-01/user-exercise.diff)是相对 HEAD 的累计差异，包含之前改名和格式整理，不将整份差异归为本轮新增。
- 学习与验证限制：给定目标、允许参考已有赋值后的相近修改及字段读取解释已核验，本项入门练习完成；不代表独立设计构造函数或完整模块维护能力。编译只支持引用与类型正确，未验证实际开始消息收发，未执行测试或联机；模式初始化仍未修复。下一项为 channel 创建与收发。

## A1.1 通道练习首轮（2026-11-11，未通过）

- 批次：`a11-channel-practice-20261111-01`；执行对话：本项目 A 当前 Phase 0 / A1.1 聊天。目标为创建容量 1 的整数通道，将同一通道保存到另一个变量，通过 `savedNumbers` 发送 90、从 `numbers` 接收，预期打印 `received: 90`。AI 仅准备临时程序外壳，用户写通道相关代码，AI 读取、执行及整理文档，没有代改用户收发代码。
- 用户[首次源码快照](evidence/a11-channel-practice-20261111-01/user-attempt.go.txt)：创建实际为 `make(chan int)`，正确赋值 `savedNumbers := numbers`，随后发送 90，接收也使用 `savedNumbers`。两个变量保存同一通道值；接收用哪个变量不会消除本次发送处的等待。练习后续将改用 `numbers` 接收，直接观察两个变量使用同一通道。
- 实际执行：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`；`go run -mod=readonly C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 输出 `fatal error: all goroutines are asleep - deadlock!`，主执行流程停在 `main.go:10` 的 channel send；Go 工具退出码 1，输出内的程序状态为 `exit status 2`。完整环境、命令及原始输出见[首次运行记录](evidence/a11-channel-practice-20261111-01/first-run.log)。
- 原因与待办：无缓冲通道发送需要接收方准备好，当前只有按顺序执行的 `main`，发送等待导致后续接收未执行。下一步由用户补容量 1 并从 `numbers` 接收，再检查实际结果；目前未通过完整收发练习。未执行项目测试或联机，这不是 R1 回归测试的修复前失败证据，也不计作业务修复贡献。

## A1.1 通道练习二轮（2026-11-11，单文件收发通过）

- 批次：`a11-channel-practice-20261111-02`；执行对话及外壳分工沿用首轮。用户按 AI 给出的两行明确修改提示，补上容量 1 并改从 `numbers` 接收；AI 未代改代码，读取实际文件后执行。归为指导下完成，不作为独立设计、业务修复或项目测试通过的证据。
- [实际源码](evidence/a11-channel-practice-20261111-02/user-buffered.go.txt)：`numbers := make(chan int, 1)` 创建通道，`savedNumbers := numbers` 保存同一通道值，`savedNumbers <- 90` 发送，`received := <-numbers` 接收并打印。容量 1 使当前顺序流程可以先暂存一个整数，再接收；两个变量访问同一条通道。
- 实际执行：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`；`go run -mod=readonly C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 输出 `received: 90`，退出码 0，见[原始运行记录](evidence/a11-channel-practice-20261111-02/buffered-run.log)。首轮失败记录保持原样。
- 学习边界：已在指导下完成正常通道保存与收发；未初始化通道的区别尚待对照，未执行项目测试或联机，R1 未修复。下一步由用户修改一行，观察声明通道变量但未给它保存通道值时的行为。

## A1.1 通道练习三轮（2026-11-11，nil 通道对照）

- 批次：`a11-channel-practice-20261111-03`；用户按明确提示，只将 `savedNumbers := numbers` 改为 `var savedNumbers chan int`，其他收发代码保持不变。AI 读取并执行，没有代改用户代码；这是指导下的负例演示，不计为个人业务修复或独立掌握。
- [实际源码](evidence/a11-channel-practice-20261111-03/user-nil-channel.go.txt)中 `numbers` 保存已创建、容量为 1 的通道，`savedNumbers` 则为 nil；`savedNumbers <- 90` 是发送整数，不会给该变量保存通道值。
- 实际执行：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`；`go run -mod=readonly C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 报告 `fatal error: all goroutines are asleep - deadlock!`，栈为 `main.go:10`、`chan send (nil chan)`，Go 工具退出 1，程序状态为 `exit status 2`。见[原始运行记录](evidence/a11-channel-practice-20261111-03/nil-channel-run.log)，正常版本与首轮输出均未覆盖。这符合本次故意移除通道赋值后的预期负例行为，未执行项目测试或联机。
- 当前学习缺口：用户提问“后面不是有赋值90吗”，需澄清保存通道值与发送消息的区别；当前不据指定编辑完成判定 channel 已掌握。理解后再恢复正常赋值并继续项目，R1 仍未修复。

## A1.1 模式通道字段赋值（2026-11-11，待完整验收）

- 批次：`a11-mode-field-edit-20261111-01`；执行对话：本项目 A 当前 Phase 0 / A1.1 聊天。AI 已示范参数与字段交接、普通通道与 nil 负例，随后给出真实代码的交接定位任务；用户自行在 `NewHub` 中补入 `setModeCh: setModeCh` 并回传截图，AI 没有提供这行完整补丁或代写业务修复。不是独立设计整个模块，完整测试断言仍待用户编写。
- 设计理由与调用链：复用已有初始化，把调用方创建、通过参数传入的模式通道值保存到 `Hub.setModeCh`，供 `ForwardSetMode` 读取。没有新增通道、抽象或业务功能。上游 Hub、消息协议及游戏循环仍为上游实现。
- 保存证据：[当前源码快照](evidence/a11-mode-field-edit-20261111-01/user-hub.go.txt)、[累计差异](evidence/a11-mode-field-edit-20261111-01/user-hub.diff)、[用户截图](evidence/a11-mode-field-edit-20261111-01/user-mode-field-edit.png)。截图副本与附件 SHA256 一致；差异相对 HEAD，包含之前参数改名与格式整理，不能全算作本轮新增。
- 实际验证：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`，HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；`git diff --check -- internal/server/hub.go internal/server/hub_test.go`、`go build -mod=readonly ./internal/server` 均无输出、退出码 0，见[原始记录](evidence/a11-mode-field-edit-20261111-01/build.log)。只验证编译，未执行项目测试或联机。
- 学习与顺序限制：用户正确定位字段原为 nil，但推断“所以一直是0”需纠正；游戏循环默认模式为 1，nil 通道导致原发送分支不可执行，现有 `default` 跳过后未收到模式更新。字段修复已经在完整回归测试运行前写入，不能宣称按时间先红后绿。保留此修改，后续在仓库外旧实现副本上运行同一完整测试做失败对照，再验证当前修复版，真实顺序另记。当前 A1.1 未完成验收。

## A1.1 select/default 演示（2026-11-11，指导下）

- 批次：`a11-select-default-20261111-01`；AI 提供完整的小示例，用户写入临时练习。用户将完整 `main` 粘贴到原 `main` 内，AI 保存[原始粘贴](evidence/a11-select-default-20261111-01/user-paste-before-wrapper-fix.go.txt)，仅删除重复的样板函数外壳并执行 gofmt；通道变量及 select 分支未改。原始重复外壳版本未执行，不编造其编译输出。
- [实际执行源码](evidence/a11-select-default-20261111-01/executed-source.go.txt)：`savedNumbers` 为 nil；select 的发送分支无法执行，default 打印“发送被跳过”，随后打印“函数继续执行”。这是指导下演示，不记为用户独立设计或独立掌握。
- 实际执行：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`；`gofmt -w C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 无输出、退出 0；`go run -mod=readonly C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 依次输出“发送被跳过”“函数继续执行”，退出 0。见[原始运行记录](evidence/a11-select-default-20261111-01/select-default-run.log)。没有消息送出，但程序正常结束；这与项目空 default 使模式发送被跳过的路径对应。
- 边界与下一步：未执行项目测试或联机，未完成 A1.1 验收。接着由用户自行将临时通道初始化改为容量 1，保留同一 select，预测并运行检查分支变化；不预先提供完整修改答案。

## A1.1 select/default 相近修改（2026-11-11）

- 批次：`a11-select-default-20261111-02`；用户得到“将 savedNumbers 初始化为容量 1 的整数通道、保留 select”的目标，可参考前例，未收到完整替换行。用户自行修改实际文件，并预测“发送成功 继续执行”；AI 读取实际文件并运行，没有修改本轮练习代码。
- [实际源码](evidence/a11-select-default-20261111-02/user-buffered-select.go.txt)使用 `savedNumbers := make(chan int, 1)`。通道有效且缓冲为空，发送可以立即进行，因此进入发送分支，随后继续执行 select 后的打印。相较前一批次 nil/default，仅改变初始化。
- 实际执行：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`；`go run -mod=readonly C:\Users\c8618\AppData\Local\Temp\a11-channel-20261111\main.go` 依次输出“发送成功”“函数继续执行”，退出码 0，见[原始记录](evidence/a11-select-default-20261111-02/buffered-select-run.log)。预测与实际分支一致，旧批次未覆盖。
- 掌握与验证边界：这个给定目标后的局部修改和输出预测已核验；分支原因由 AI 补充解释，尚不能据此认定独立掌握并发或测试设计。没有执行项目测试或联机。下一步由用户编写模式回归测试的消息内容断言。

## A1.1 回归断言首次执行（2026-11-11，格式检查失败）

- 批次：`a11-regression-draft-20261111-01`；用户根据整数检查示例和预期字段目标，自行将输入改为玩家 1、保存消息至 `msg`，写入 `msg.PlayerID != 1` 与 `msg.Mode != protocol.ModeLastStanding` 两条判断。AI 只读检查和执行，没有代写或改动这些断言。见[实际测试源码](evidence/a11-regression-draft-20261111-01/user-regression-draft.go.txt)与[当前累计差异](evidence/a11-regression-draft-20261111-01/working-changes.diff)。
- 实际命令：Windows amd64、Go 1.27.3，目录 `C:\Dev\Go\A`，HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；`go test -mod=readonly ./internal/server -run '^TestHubForwardSetModeSendsMessage$' -count=1` 退出 1。见[原始输出](evidence/a11-regression-draft-20261111-01/targeted-test.log)。
- 原因：两条判断本身符合预期，但第 212 行 `t.Error` 的提示含 `%d`，第 215 行 `t.Fatal` 的提示含 `%v`；Go vet 分别报告 `Error call has possible Printf formatting directive %d`、`Fatal call has possible Printf formatting directive %v`，包报告 `[build failed]`，未执行断言。
- 边界与下一步：此时工作文件已有用户字段修复，这次失败来自测试提示方法，不能作为原模式缺陷的失败对照。由用户把两处报告方法改为 Errorf / Fatalf，判断条件与无消息 default 保持原样；之后再运行旧实现及当前修复版的同一测试。全库测试与修复后联机未执行。

## A1.1 正式回归与全库验证（2026-11-11，联机待验收）

- 批次：`a11-regression-20261111-01`；用户编写字段赋值、非零玩家编号输入及两条字段判断，AI 解释并指导格式化报告方法、准备版本副本、辅助执行和记录。没有由 AI 代写核心修复或断言。本次自动验证完成，A1.1 尚待本人结果解释及修复后联机复验。
- 对照方法：以原始提交 `488dd313bc13de387d12a24f5d05a5c28bed5899` 执行 git archive，解至仓库外临时目录，只把其中的 hub_test.go 替换为当前完整测试。旧版与当前版使用相同测试，SHA256 为 `E99A69D8BBC9C015F03737299787255AAF997A0BCADFE64EFE5271460CB9C02C`；输入为 PlayerID 1、ModeLastStanding，检查收到消息及两个字段。见[旧 Hub](evidence/a11-regression-20261111-01/old-hub.go.txt)、[当前 Hub](evidence/a11-regression-20261111-01/current-hub.go.txt)、[共用测试](evidence/a11-regression-20261111-01/shared-hub-test.go.txt)和[累计差异](evidence/a11-regression-20261111-01/working-changes.diff)。
- 环境与顺序：[环境及副本校验记录](evidence/a11-regression-20261111-01/environment.log)包含 Go 1.27.3 / Windows amd64、原始提交、临时位置和哈希。GOWORK=off 只设于验证进程；依赖版本未改。用户先写入修复，后完成测试，因此这里是旧/新实现版本对照，不冒充按时间先红后绿。先前 vet 失败单独保留，不当作旧缺陷失败。

| 验证 | 实际命令与结果 | 原始记录 |
| --- | --- | --- |
| 旧实现定向测试 | 在临时旧实现运行 `go test -mod=readonly ./internal/server -run '^TestHubForwardSetModeSendsMessage$' -count=1 -v`；第 219 行报 `mode channel received no message`，退出 1。测试确实执行，在无消息 default 失败。 | [旧版输出](evidence/a11-regression-20261111-01/old-targeted-test.log) |
| 当前修复版定向测试 | 在 A 运行相同命令；该测试通过，退出 0，说明本次收到的 PlayerID / Mode 符合预期。 | [当前版输出](evidence/a11-regression-20261111-01/current-targeted-test.log) |
| 相关服务器包 | `go test -mod=readonly ./internal/server -count=1 -v`；5 个测试通过，退出 0。 | [包测试](evidence/a11-regression-20261111-01/server-package-tests.log) |
| 全库测试 | `go test -mod=readonly ./... -count=1 -v`；4 个有测试包通过，37 个顶层测试 PASS、0 SKIP、0 FAIL，另 4 包无测试文件，退出 0。 | [全库测试](evidence/a11-regression-20261111-01/full-tests.log) |
| 全库构建 | `go build -mod=readonly ./...`；无输出、退出 0。 | [构建](evidence/a11-regression-20261111-01/full-build.log) |
| 范围检查 | 依赖及 LICENSE 差异为空；受跟踪修改只有 hub.go / hub_test.go；空白检查退出 0，有 Git LF/CRLF 提示。 | [范围审查](evidence/a11-regression-20261111-01/scope-audit.log) |

- 格式与副本准备：初始 gofmt 检查列出两文件，先暂停并保存[格式差异](evidence/a11-regression-20261111-01/gofmt-before.diff)及格式化前源码。AI 执行 gofmt 后逐字对比移除 CR 的文本，确认仅 CRLF 转为 LF，核心语句未改。原始 archive 文件受 core.autocrlf=true 影响，原始字节哈希不同；经仓库换行规范化后的 Hub blob 与原始提交完全一致，正文对比也一致。两次准备检查暂停均发生在正式定向测试开始前，未当作测试失败。
- 局限与待办：单元回归直接验证 Hub 转发，不包含真实 WebSocket、房主权限或游戏循环采用模式的全过程。修复后双客户端联机尚未进行；无 Race、压力或性能结论。测试时间仅保留为命令原始输出，不作为性能改进数据。本轮结束前[8080 点查](evidence/a11-regression-20261111-01/port-8080.log)无监听，只代表该时刻；AI 未启动交互服务。临时旧副本不保证长期保留，关键源码与日志已归档。

## A1.1 本人修复后联机（2026-11-11，工程验证完成）

- 批次：`a11-snake-after-fix-20261111-01`；用户亲手操作服务端与两个客户端，AI 提供目录、命令、菜单、开局及退出指引，核对附件和用户反馈。指引命令为在 `C:\Dev\Go\A` 分别运行 `go run ./cmd/server`、`go run ./cmd/client -name alice`、`go run ./cmd/client -name bob`，两端选择 Last Standing，alice 开局，客户端 Q 退出、服务端 Ctrl+C。截图未完整覆盖输入命令，未采集手工命令退出码。
- [服务端原始截图](evidence/a11-snake-after-fix-20261111-01/server-mode-2-and-clean-shutdown.png)从用户附件复制，SHA256 与源文件一致。下表时间按截图原样记录，未另核对终端时区，不作为性能数据。

| 观察点 | 截图证据 |
| --- | --- |
| 加入与模式更新 | 06:10:44 alice 注册为 id=0、成为房主；06:10:45 `game loop: mode set to 2`；06:11:09 bob 注册为 id=1、玩家数为 2。 |
| 双人开局 | 06:22:53 `game loop: starting — mode=2 players=2`，随后 `ticking level 1 mode 2`。修复前实际模式为 1，本次实际采用 Last Standing 的模式 2。 |
| 连接退出 | 06:23:25 客户端 0 发出 leave，注销至 total=1；06:23:26 客户端 1 发出 leave，注销至 total=0；两连接 readPump / writePump 退出。 |
| 服务端停机 | 06:23:33 出现 `shutting down snake server...`、`snake server stopped cleanly`，返回 `PS C:\Dev\Go\A>`。 |

- 客户端侧补充：针对“两端是否都显示过棋盘，并各按 Q 返回各自 PowerShell 提示符”的明确问题，用户回答“是的”。据此记录两端棋盘与 Q 退出已由用户确认，AI 未直接观察客户端画面；未单独采集客户端截图、退出码或完整方向键操控证据，不据此宣称完整游戏规则通过。
- [清理及源码核对](evidence/a11-snake-after-fix-20261111-01/cleanup-and-source-check.log)：AI 在用户停机后实际检查 8080 无监听，hub.go / hub_test.go 与已通过自动测试的快照在换行规范化后完全一致，因此未重复测试。没有 AI 强制终止进程、提交、推送、部署或依赖/许可证变更。
- 当时结论：字段修复、回归测试、全库构建与本轮有限联机冒烟均有真实证据，工程验证完成。用户复盘指出初始化问题，但仍表示不知道为何没卡住，且将发送语句与读消息混淆；之后补讲的学习收尾见下一节，不把工程通过写成已经独立掌握并发模块。

## A1.1 学习收尾与交付（2026-11-11）

- 执行对话：本项目 A 的开发与教学聊天。延续上节实际缺口，AI 区分 `通道 <- 消息` 的发送和 `消息 := <-通道` 的接收，并说明 select 选择能立即进行的通信；没有可进行的通信时，有 default 才走 default。旧 Hub 字段为 nil，发送分支不能执行，空 default 使函数继续，不代表消息送达。
- 用户先用 C# switch/case 作类比，对已满通道的发送条件判断仍有混淆；AI 补充 switch 匹配值或模式、select 判断通信是否可立即进行。最后对“容量 2 的缓冲通道已有 1 条消息，再执行带 default 的发送”这一相近问题，用户回答“有空位，能发送”，该局部边界理由正确。此项是聊天中的口头预测，没有新增程序运行或测试数据。
- 实际学习结果：参数改名和字段交接、容量 1 通道初始化、真实字段修复、玩家 1 输入与两字段断言均有用户保存的实际修改和 AI 检查/执行证据。其中相近修改由用户依据目标或示例完成，格式化报告方法经指导修正。记录为指导下完成、局部修改已验证；不据最后一道预测宣称全面独立掌握 channel、并发或完整测试设计。发送/接收、通道值与消息内容、完整链路解释仍需后续实践巩固。
- 交付范围：用户修复位于 hub.go，回归位于 hub_test.go；沿用上游架构，恢复已有通道的保存，无新增抽象。模式选择 → readPump → ForwardSetMode → 模式通道 → Snake 游戏循环的本次链路已由单元回归及有限联机证据支持。没有新增完整规则、Race、性能或跨机器结论。
- A1.1 于 2026-11-11 收尾并停止；工程验证批次及旧/新版本对照保留如上，不把先写修复的实际过程改写成按时间先红后绿。联机后未修改业务代码、未重复测试，本次只更新本地状态、证据与阶段约定。没有依赖/许可证/共享文件修改、提交、推送或部署；下一任务须由用户明确开启。

## A1.2 任务开启与源码调查（2026-11-11，尚无运行结果）

- 执行对话：本项目 A 接续聊天。用户明确“开始吧”，批准上一轮提出的扫雷光标消息身份归属修复，现记为 A1.2；核心修复和代表性断言继续由用户在指导下编写，AI 辅助调查、讲解、验证和本地文档。尚无 A1.2 用户代码贡献或测试通过结果。
- 只读实查：工作目录与 Git 根目录均为 `C:\Dev\Go\A`，远程来源为 `https://github.com/siyad01/multiplayer-terminal-games.git`，HEAD 为 `488dd313bc13de387d12a24f5d05a5c28bed5899`；开启前受跟踪修改仅为 A1.1 的 hub.go / hub_test.go，保留未提交成果。未连接远程或操作凭据。
- 静态路径：`protocol.EncodeCursor` 编码 `[type][playerID][x][y]`；扫雷入口 `MsgMineCursor` 将消息的 `data[1]` 保存到 `MineInputMsg.PlayerID`，游戏循环 `ActionCursor` 将其传给 `Game.MoveCursor` 更新对应光标并广播。翻格、插旗的相邻入口使用 `c.ID()`。这支持定位身份来源差异，尚未运行伪造消息或证明完整网络行为。
- 现有测试主要覆盖游戏方法及协议编码，没有发现直接覆盖扫雷 WebSocket 身份选择的测试。后续回归须覆盖真实入口的身份归属，不以单独测试 MoveCursor 替代；最小测试接缝仍在调查，尚未提取处理函数。
- 本次仅更新本地 AGENTS、CURRENT_STATE 和本文，说明授权与当前步骤；未修改业务源码或测试、未运行测试/构建/联机、未安装组件、未提交/推送/部署。A1.1 的既有结果不算 A1.2 新验证。

## A1.2 回归入口与测试外壳准备（2026-11-11，待用户身份断言）

- 批次：`a12-cursor-scaffold-20261111-01`；执行对话为本项目 A 接续聊天。用户先根据入口代码回答伪造编号会移动“bob 的”光标，局部预测正确。AI 随后准备下述测试接缝和网络样板；尚未发生 A1.2 用户实际代码编辑，不将预测记为独立修复或独立测试设计。
- 改动与原因：AI 将 `cmd/mineserver/main.go` 的现有 `/ws` 闭包提取为 `mineWSHandler`，以便测试直接挂载生产同一入口并使用临时端口，无需启动带信号处理的 main 或固定 8081 服务。没有新增接口或通用解析层。光标入口仍取 `data[1]`；在 gofmt 前，提取前后的处理函数体逐字比较为 True，见 [提取核对](evidence/a12-cursor-scaffold-20261111-01/extraction-check.log)、[原始源码](evidence/a12-cursor-scaffold-20261111-01/main-before-extraction.go.txt)、[提取后源码](evidence/a12-cursor-scaffold-20261111-01/main-after-extraction.go.txt)及[差异](evidence/a12-cursor-scaffold-20261111-01/main-extraction.diff)。本轮其他变化仅为所改文件 gofmt，不计为用户核心修复。
- AI 辅助测试结构：新增 `cmd/mineserver/main_test.go`，真实 Hub 注册两名 WebSocket 客户端，通过 Welcome 与 join 回复同步，从 Alice 连接发出填有 Bob 编号的 Cursor；从实际 `mineCh` 等待消息，并检查动作及坐标。接收使用有界等待，无 Sleep。核心身份断言留有确定失败的 `t.Fatal("TODO: 请补充连接身份断言")`，等待用户替换；尚未执行它，不把占位失败当作原缺陷失败。见[待用户编辑的外壳快照](evidence/a12-cursor-scaffold-20261111-01/scaffold-before-learner-assertion.go.txt)。
- 实际命令及结果：Windows amd64、Go 1.27.3、工作目录 `C:\Dev\Go\A`、HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`。`gofmt -w cmd/mineserver/main.go cmd/mineserver/main_test.go` 无输出、退出 0，见[格式记录](evidence/a12-cursor-scaffold-20261111-01/format.log)；进程级 GOWORK=off 下执行 `go test -mod=readonly ./cmd/mineserver -run '^$' -count=1`，输出 `ok ... [no tests to run]`、退出 0，见[编译检查](evidence/a12-cursor-scaffold-20261111-01/compile-only.log)。仅编译和默认静态检查，不是用例执行或回归通过。未运行全库测试、完整构建或交互联机。
- 范围与清理限制：保留 A1.1 的两文件内容，源码哈希及范围核对另存本批次；无依赖、许可证、共享文件、提交、推送、部署或系统配置修改。未来运行外壳时会显式关闭 WebSocket、等待 leave 注销并关闭临时 HTTP 服务；现有 Hub.Run 无停止接口，仍会有阻塞 Hub goroutine 留到测试进程结束，不宣称完整 goroutine 回收，也不扩大为 Hub 生命周期改造。这个外壳只验证 WebSocket 至游戏输入通道，不替代完整游戏光标广播与绘制的联机验收。
- 下一步：用户在 `TestMineCursorUsesConnectionIdentity` 的 TODO 位置写入 `msg.PlayerID` 与 Alice 连接编号的判断；AI 读取实际保存文件，运行同一用例并保存修复前真实身份不符的失败，再指导用户修复入口。

## A1.2 用户身份断言与修复前失败（2026-11-11）

- 批次：`a12-cursor-before-fix-20261111-01`；执行对话为本项目 A 接续聊天。用户按给出的三行示例，在 `TestMineCursorUsesConnectionIdentity` 末尾保存 `if msg.PlayerID != aliceID` 和 `t.Fatalf`，自行填写提示文本。AI 读取实际文件后执行，未代写或改动该断言。记录为指导下完成，不记为从空白独立搭建测试。
- 运行前快照：[尚未修复的入口](evidence/a12-cursor-before-fix-20261111-01/main-before-fix.go.txt)、[用户保存的实际断言](evidence/a12-cursor-before-fix-20261111-01/user-identity-assertion.go.txt)、[源码哈希](evidence/a12-cursor-before-fix-20261111-01/source-hashes-before.json)。入口仍为 `PlayerID: data[1]`，没有先行修复；两源码在执行前后哈希一致。`gofmt -l cmd/mineserver/main.go cmd/mineserver/main_test.go` 无输出、退出 0，没有运行 gofmt 写入，见[格式检查](evidence/a12-cursor-before-fix-20261111-01/format-check.log)。
- 实际环境与命令：Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；GOWORK=off 仅设于验证进程，结束后还原。执行 `go test -mod=readonly ./cmd/mineserver -run '^TestMineCursorUsesConnectionIdentity$' -count=1 -v -timeout=20s`，退出 1。完整输出见[修复前测试日志](evidence/a12-cursor-before-fix-20261111-01/targeted-test.log)。
- 实际结果：用例进入 RUN，Alice 注册为 0、Bob 为 1；Alice 连接发送填有 Bob 编号的 Cursor，实际 `mineCh` 收到 `{PlayerID:1 Action:2 X:3 Y:4}`。动作和坐标检查通过，随后用户在第 112 行的身份断言报告“期望0,实际1”，该用例 FAIL。没有编译/vet、网络超时或 TODO 占位失败，可以作为身份归属缺陷的真实修复前证据；原始耗时不作为性能数据。
- 结束与边界：测试清理日志显示两客户端注销至 total=0，没有清理超时；随后测试进程退出。未核验所有 pump goroutine 已结束，Hub.Run 的既有限制仍按上一节说明。不包含完整游戏循环、广播绘制的交互验收；未运行全库测试、完整构建或 Race。本轮未改业务源码或测试，仅归档证据并更新本地状态；下一步由用户修改光标入口的 PlayerID 来源，再检查同一用例。

## A1.2 用户修复与自动验证（2026-11-11，联机待验收）

- 批次：`a12-cursor-after-fix-20261111-01`；执行对话为本项目 A 接续聊天。用户按明确的一行替换指引，将 `mineWSHandler` 的 `MsgMineCursor` 分支改为 `PlayerID: c.ID()` 并保存；AI 读取实际文件后验证，未代写此修复。核心断言和修复均为指导下完成，处理函数提取及网络测试外壳为 AI 辅助，不宣称用户独立设计模块。
- 最小修复与顺序：消息格式、动作和坐标处理不变，只把身份来源改为当前连接。与修复前快照按 CRLF/LF 统一后比较，正文只此一行变化；测试文件字节 SHA256 保持 `C7561F0E6000FBF5A13FB163EDF663E8FEA1CDD60E1E8067629A9313DE3E1FDE`。见[同源核对](evidence/a12-cursor-after-fix-20261111-01/before-after-source-check.log)、[用户修复源码](evidence/a12-cursor-after-fix-20261111-01/user-fixed-main.go.txt)及[相同回归测试](evidence/a12-cursor-after-fix-20261111-01/same-regression-test.go.txt)。实际顺序是先运行用户断言取得真实身份失败，再由用户保存修复并运行相同测试。
- 环境：Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `488dd313bc13de387d12a24f5d05a5c28bed5899`；验证进程 GOWORK=off，退出后还原。gofmt 只读检查无输出、退出 0，未重写用户代码，见[格式检查](evidence/a12-cursor-after-fix-20261111-01/format-check.log)。

| 验证 | 实际命令与结果 | 原始记录 |
| --- | --- | --- |
| 同一身份回归 | `go test -mod=readonly ./cmd/mineserver -run '^TestMineCursorUsesConnectionIdentity$' -count=1 -v -timeout=20s`；Alice=0 的连接仍发送 Bob=1 的载荷，实际收到 `{PlayerID:0 Action:2 X:3 Y:4}`，用例 PASS、退出 0。两连接注销至 total=0。 | [定向测试](evidence/a12-cursor-after-fix-20261111-01/targeted-test.log) |
| 全库测试 | `go test -mod=readonly ./... -count=1 -v -timeout=60s`；38 个顶层 PASS、0 SKIP、0 FAIL，5 个有测试包、3 个无测试包，退出 0。 | [完整输出](evidence/a12-cursor-after-fix-20261111-01/full-tests.log)、[计数](evidence/a12-cursor-after-fix-20261111-01/full-test-summary.json) |
| 全库构建 | `go build -mod=readonly ./...`；无输出、退出 0，未在仓库生成可执行文件。 | [构建](evidence/a12-cursor-after-fix-20261111-01/full-build.log) |
| 源码与范围 | 两个本任务源码在测试与构建前后哈希一致，A1.1 两源码与本任务开启前哈希一致；依赖/许可证差异为空，空白检查退出 0，有 Git LF/CRLF 提示。 | [范围核对](evidence/a12-cursor-after-fix-20261111-01/scope-audit.log) |

- 局限及下一步：回归真实覆盖 WebSocket 至游戏输入通道的身份归属，未启动完整游戏循环；正常双客户端光标广播和绘制仍待用户操作验收。8081 只读点查无监听，仅代表当时，见[端口记录](evidence/a12-cursor-after-fix-20261111-01/port-8081.log)；AI 未启动长期服务。无 Race、压力、性能或完整规则结论，原始测试耗时不作性能数据。既有 Hub.Run 无停止接口的测试进程内局限不变；无依赖/许可证/共享文件改动、提交、推送、部署或系统安装。

## 个人仓库建立与远程配置（2026-11-11）

- 执行：A 项目本聊天；用户要求在已连接的本人 GitHub 下建立项目仓库，并询问本地 Git 是否存在。AI 负责本次仓库配置，未接管 A1.2 核心修改。该操作不计作个人业务功能成果。
- 本地核对：`git rev-parse --is-inside-work-tree` 返回 true，Git 目录为 `C:/Dev/Go/A/.git`，HEAD 为 `488dd313bc13de387d12a24f5d05a5c28bed5899`。这是已有克隆，无需重新执行 git init；用户的修改仍未提交。
- GitHub 连接工具没有创建 Fork 的接口，改用本机 GitHub CLI。`gh auth status --active --hostname github.com --json hosts --jq ...` 确认账号 LearnedYet、state=success；未查看或修改真实凭据。创建前本人同名仓库查询为 HTTP 404。
- 首次 `gh repo fork siyad01/multiplayer-terminal-games --clone=false --remote=false` 因 CLI 不允许在指定仓库参数时使用 --remote 而退出 1，未创建；去掉该参数后 `gh repo fork siyad01/multiplayer-terminal-games --clone=false` 退出 0，创建 [个人 Fork](https://github.com/LearnedYet/multiplayer-terminal-games)。API 随后确认 fork=true、private=false、parent=siyad01/multiplayer-terminal-games，账号具有管理及推送权限；没有通过实际 push 验证传输。
- 本地依次执行 `git remote rename origin upstream`、`git remote add origin https://github.com/LearnedYet/multiplayer-terminal-games.git`、`git config --local remote.pushDefault origin`，均退出 0。随后核对两个远程的 fetch/push 地址及默认推送远程，HEAD 未变。未克隆第二份仓库、未提交、未上传本地代码、文档或截图。
- 独立批次为 `git-personal-fork-20261111-01`，[操作及结果](evidence/git-personal-fork-20261111-01/setup.log)保留真实成功和失败。只更新本地 Git 配置及本节/当前状态，没有修改许可证、依赖、共享文件或运行程序测试；A1.2 的开发与验收保持其自身记录。

## A1.2 本人双客户端光标观察及停机（2026-11-11，工程验证完成）

- 批次：`a12-cursor-manual-20261111-01`；执行对话为本项目 A 接续聊天。用户按上一轮服务端、Alice、Bob 三窗口和 Cooperative / Easy 指引操作，回传两张双端棋盘截图及一张服务端停机截图。AI 核对画面、渲染代码及当前源码，并归档附件；没有 AI 操控终端或新增运行测试。
- 原始附件：[第一组画面](evidence/a12-cursor-manual-20261111-01/01-two-clients-after-alice-move.png)、[第二组画面](evidence/a12-cursor-manual-20261111-01/02-two-clients-after-bob-move.png)、[服务端停机](evidence/a12-cursor-manual-20261111-01/03-server-clean-shutdown.png)。从用户给出的临时 PNG 原样复制，三份源文件及副本 SHA256 均一致，见[附件清单](evidence/a12-cursor-manual-20261111-01/attachment-manifest.json)。没有覆盖旧批次。
- 棋盘观察：两端显示 MINESWEEPER – COOPERATIVE，均列出 alice / bob、分数 0。第一组中 Bob 端显示 Alice 在第一行第三格；第二组中两端最终对应 Alice (2,0)、Bob (0,2)（坐标从 0 开始）。这些静态画面支持正常光标位置到达另一端的有限观察，不等同完整按键录像或逐帧一致性证明。`render.go` 的 DrawBoard / drawCell 会将本人光标用 Navy 背景、对方用玩家颜色显示，颜色差异不作为身份错误；第一组 Bob 端还可见本人本地位置与服务端初始位置的叠加，不据此扩大为本任务新的修复范围。
- 服务端截图直接显示实际命令 `go run -mod=readonly ./cmd/mineserver`。20:02:53 alice 注册 id=0，20:03:00 bob 注册 id=1；20:03:07 `starting mode=1 diff=easy players=2`，随后 `board 9x9 mines=10`；20:04:46 / 20:04:47 两人注销，最终 total=0；20:04:58 stopped cleanly 并返回 `PS C:\Dev\Go\A>`。时间按截图原样保留，未另核对终端时区，不作性能数据。客户端启动命令、具体移动/退出按键、退出后各自提示符和进程退出码未单独采集；不补写为直接观察。
- 停机后 AI 实际执行只读端口及哈希检查：8081 无监听，main.go / main_test.go 与已通过自动验证的源码 SHA256 一致，见[清理及源码核对](evidence/a12-cursor-manual-20261111-01/cleanup-and-source-check.log)。因此未重复测试。没有强制终止进程、源码改动、提交、推送、部署或安装。
- 结论及边界：核心身份回归、全库测试/构建以及本次正常光标位置同步与停机的有限证据已齐，工程验证完成。当前仍待用户简短复盘身份与坐标来源，核心代码和断言维持“指导下完成”归属，不升级为独立网络模块设计。没有完整游戏规则、Race、压力、性能或跨机器结论。

## A1.2 学习复盘与交付（2026-11-11）

- 用户用自己的话说明：玩家编号由服务器分配，若采用客户端消息中的编号，修改客户端即可冒用他人的编号移动光标；X、Y 表示客户端提交的光标坐标，因此从消息取得。该解释正确区分了连接身份和操作参数。AI 补充坐标仍是需要按游戏规则校验的外部输入，不将身份修复扩大为已完成坐标校验或完整安全治理。
- 实际贡献归属：用户按明确示例在实际文件保存代表性身份断言和 `PlayerID: c.ID()` 单行修复，亲手运行双客户端并提交截图。AI 提取现有处理函数、搭建网络测试外壳、提供代码指导、执行验证和整理证据。复盘是口头解释，不是新增独立编码或运行数据，未据此认定独立设计整个网络模块或测试框架。
- 交付依据：同测试按实际时间取得修复前失败、修复后通过，全库 38 个顶层测试通过及构建成功；有限双端正常光标位置同步和正常服务端停机有原始截图，停机后 8081 无监听。相关批次见上文。收尾时再次只读核对 main.go / main_test.go 与自动验证快照 SHA256 一致，本轮未重测或改业务代码。
- A1.2 于 2026-11-11 完成并停止。本轮仅更新本地 AGENTS、CURRENT_STATE 和本文，保留 BASELINE 及原始证据；没有提交、推送、部署、依赖/许可证/共享文件修改或系统安装。后续任务尚未授权。完整规则、Race、性能及跨机器验证等限制继续保留。

## A1 本地提交及 A2 开启（2026-11-11）

- 用户明确要求“提交，然后开启下一阶段”。AI 核对现有暂存区、最新交付文档及源码：范围为 92 个项目约定、文档/证据及四个业务/测试文件；无 LICENSE、go.mod、go.sum 改动。main.go / main_test.go 与已通过 38 项测试的快照一致，hub.go / hub_test.go 与此前记录一致，因此未重复运行程序测试。
- 更新暂存区中的三份最新交接文档后，执行 `git commit -m "fix: preserve snake mode channel and bind mine cursor to connection identity"` 成功，提交为 `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`。92 files changed，5694 insertions、54 deletions；提交后、A2 文档开始前 `git status --porcelain` 为空。见[本地提交回执](evidence/a2-start-20261111-01/local-commit-and-phase-start.log)。未推送，未改远程或凭据；Git 创建提交不改变用户指导下编写核心代码的贡献归属。
- 随后开启 A2 单对局生命周期的读码/规则确认。实际读取 minesweeper Loop/runOnce 和 Game 的玩家保存/删除路径，确认现有离开处理会删除含本局分数的玩家对象、Cooperative 允许中途加入、全员离开未见取消结果。只属于静态调查，没有新的运行或功能结论。
- 已在 CURRENT_STATE 列出目标、用户核心作者、AI 范围、验收方向及固定开局名单的推荐方案/允许中途加入的备选；参与、取消、再次开局细节仍待确认，未编写 A2 核心代码或测试。此提交之后的 A2 启动文档属于新的工作区修改，尚未提交；未改共享文档或原始 BASELINE。

## A2 首版规则确认及既有改动核对（2026-11-12）

- 用户对上一轮首版方案回答“同意”，确认固定开局名单、离开保留成绩、全员离开取消、结果固定及再开局新编号；不再将这些业务规则标为待批准。核心作者与 AI 辅助范围沿用既有学习分工。
- 本聊天开始本轮时，git status 已显示 game.go 和 game_test.go 有未提交改动。只读 diff 看到新增 MatchId、Left、活动检查、全员离开检查，以及 TestPlayerLeaveRetainsScore。用户随后说明这是其在 Gemini 指导下编写的代码，尚未提交，待完成后交本聊天审查。按用户在 Gemini 指导下完成记录，不记为本聊天实现或用户独立设计；没有覆盖、撤回或重写这些文件。
- 首次静态审查发现当时的分数保留判断仅验证非零，没有严格比较离开前后值；退出后插旗只看玩家 Flags，没有接收并检查返回的 changes。只是未完成快照的审查记录，没有新增运行结论；用户仍在 Gemini 指导下实现，完成交接后须按届时实际文件重新核对。本聊天按用户要求等待交接，不接手源码修改或提前运行测试。
- 本轮只更新本地约定、状态及本文，未执行 Go 测试/构建、提交、推送或部署；未修改核心代码、依赖、许可证、共享文件或其他项目。之前 38 PASS 仅适用于 A1 已验证快照，不适用于这批尚未验证的 A2 修改。

## A2 接回带学与首轮验证（2026-11-12）

- 用户说明此前因本聊天额度不足由 Gemini 暂时指导，现明确要求本聊天接手。接续范围是已批准的 A2 单房间 Cooperative / Easy 生命周期；核心代码及代表性断言仍由用户编写，AI 负责讲解、审查、验证及本地文档。当前 game.go / game_test.go 改动保留“用户在 Gemini 指导下编写”归属，本轮 AI 没有修改源码或断言，没有重做既有实现。
- 只读核对：RemovePlayer 标记 Left、保留玩家对象并删除光标；Reveal / Flag / MoveCursor 检查活动玩家；AllPlayersLeft 提供全员离开判断。调用链仍为 loop.go leave → Game.RemovePlayer，尚未接入 AllPlayersLeft；游戏中 join 仍会 AddPlayer，不符合已确认的固定名单规则，取消和固定结束结果尚未完成。
- 批次为 `a2-handoff-check-20261112-01`。验证前保存[业务源码](evidence/a2-handoff-check-20261112-01/game.go.txt)、[测试源码](evidence/a2-handoff-check-20261112-01/game_test.go.txt)及[相对 HEAD 的源码差异](evidence/a2-handoff-check-20261112-01/source-diff.patch)。HEAD 为 `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，GOWORK=off 只设于验证进程并还原，见[环境](evidence/a2-handoff-check-20261112-01/environment.log)。两源码运行前后 SHA256 一致，见本批次 before / after 哈希文件。
- 实际执行：`go test -mod=readonly ./internal/minesweeper -run '^TestPlayerLeaveRetainsScore$' -count=1 -v -timeout=30s`：1 项 PASS、退出 0，见[定向日志](evidence/a2-handoff-check-20261112-01/targeted-test.log)；`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：10 项顶层 PASS、0 SKIP、0 FAIL、退出 0，见[扫雷包日志](evidence/a2-handoff-check-20261112-01/package-tests.log)及[计数](evidence/a2-handoff-check-20261112-01/test-summary.json)。本轮是已有实现的首次本聊天验证，没有取得修改前失败对照；未核验或补写 Gemini 期间的执行结果。
- 审查与下一学习点：TestPlayerLeaveRetainsScore 已保存 scoreBeforeLeave，但离开后只查非零，不能识别分数从 5 错变成 1；下一步由用户改为严格比较离开前后值。退出后插旗断言仅看 Flags，仍需后续补强。编号目前依赖毫秒时间及游戏 rng，只有非空断言，连续两局唯一性未验证，后续编号步骤再处理。
- 结论仅为当前早期模型用例及相关包通过，不代表 A2 完成或用户独立掌握。未执行全库测试、构建、Race、联机或性能验证；无源码/依赖/许可证/共享文件改动，无提交、推送、部署或系统安装。更新本地 AGENTS、CURRENT_STATE 和本文，继续等待用户保存下一小步断言。

## A2 用户补强离开前后成绩断言（2026-11-12）

- 用户按本聊天给出的比较示例，实际保存 `TestPlayerLeaveRetainsScore` 中 `g.Players[0].Reveals != scoreBeforeLeave`，并自行调整提示文本与参数顺序为离开后、离开前。AI 读取实际文件，确认比较与报告对应正确；未代写实际断言或修改业务实现。记录为指导下完成，不认定独立设计整个测试。
- 批次 `a2-score-assertion-20261112-01`：保留[用户保存的测试](evidence/a2-score-assertion-20261112-01/user-game-test.go.txt)、环境及源码前后哈希。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，GOWORK=off 仅设于验证进程并还原；执行 `go test -mod=readonly ./internal/minesweeper -run '^TestPlayerLeaveRetainsScore$' -count=1 -v -timeout=30s`，1 PASS、退出 0，见[定向日志](evidence/a2-score-assertion-20261112-01/targeted-test.log)。源码运行前后 SHA256 一致。
- 变化仅为保留成绩的断言从“非零”改成“与离开前完全相同”，可以识别分数变少但仍非零的错误；本轮没有执行注入错误的失败对照，也没有重新运行全包、全库、构建或联机，不能将本次通过升级为 A2 完成。
- 下一小步是将退出后的插旗场景单独测试：新棋盘未翻开，确保目标为隐藏格；检查 Flag 返回的变化列表及格子状态。当前只是指导，用户尚未保存该小步。AI 本轮只验证、记录及提供代码指导，无提交、推送、依赖/许可证/共享文件修改。

## A2 带学提示强度调整（2026-11-12）

- 用户报告“测试跑通了”，并明确要求后续业务改动和核心断言只提供思路提示与 TODO 骨架，让其自行敲代码，不直接贴完整答案。已记入本项目 AGENTS 和 CURRENT_STATE，后续按此执行；未修改共享学习规则或全局记忆文件。
- 本轮只读核对实际 game_test.go：严格成绩比较仍在，原 TODO(4) 插旗块仍在，未找到 TestLeftPlayerCannotFlag；其 SHA256 与上一批 a2-score-assertion-20261112-01 已验证快照一致。因此只将通过消息记为用户报告，不认定尚未保存的新插旗测试通过；本轮未重复运行 Go 测试或改业务/测试源码。

## A2 用户保存独立插旗测试并验证（2026-11-12）

- 用户说明已保存后，AI 读取实际 game_test.go，确认 TestLeftPlayerCannotFlag 已存在，原成绩测试内的插旗执行及 Flags 检查已删除（原 TODO(4) 注释尚保留）。新用例在未翻开的棋盘中让 Alice 离开，再尝试对隐藏的 (0,0) 插旗，检查返回变化列表为空及目标仍为 StateHidden。调用链为测试 → Game.RemovePlayer 标记 Left → Game.Flag 活动检查 → 返回无变化。该用例是用户根据本聊天此前完整示例保存，记为指导下完成，不视为根据 TODO 独立填出的断言；后续遵守用户刚要求的思路与骨架方式。
- 独立批次 `a2-left-flag-20261112-01` 保留[用户测试快照](evidence/a2-left-flag-20261112-01/user-game-test.go.txt)、[环境](evidence/a2-left-flag-20261112-01/environment.log)、源码运行前后 SHA256 及[计数](evidence/a2-left-flag-20261112-01/test-summary.json)。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`；GOWORK=off 仅设于验证进程并还原。实际执行 `go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`，11 个顶层 PASS、0 SKIP、0 FAIL、退出 0，见[完整包测试日志](evidence/a2-left-flag-20261112-01/package-tests.log)。源码前后哈希一致。
- 本轮 AI 未修改业务或测试源码，仅验证及更新本地状态和证据。没有注入缺陷的失败对照；未执行全库、构建、Race、联机或性能验证，A2 整体仍未完成。下一步只提供 TODO 骨架，让用户编写离开后光标不会重建的相似用例；不扩展阶段、提交或推送。

## A2 用户按 TODO 骨架补齐退出后光标测试（2026-11-12）

- 用户回传[编辑及运行截图](evidence/a2-left-cursor-20261112-01/user-cursor-test.png)，AI 随后读取实际保存的 game_test.go。TestLeftPlayerCannotMoveCursor 按给定业务思路与 TODO 骨架，由用户补齐游戏创建、加入与离开、移动请求及三次光标数量判断；每次动作后重新调用 AllCursors，未沿用旧查询结果，并删除占位 t.Fatal。AI 没有预先提供该用例的完整动作/断言答案，也未写入核心代码。可记录为在提示下自行完成这个相似小变化，不能据此认定独立设计测试体系。
- 关键调用链：AddPlayer 建立光标 → AllCursors 查询有一个 → RemovePlayer 标记离开并删除光标 → 查询为空 → MoveCursor 对离开者返回 → 再查询仍为空。单玩家场景使数量检查能识别删除失败或错误重建；该用例不验证 Bob 的光标或全局游戏结束后的行为。
- 独立批次 `a2-left-cursor-20261112-01` 保留[用户保存的测试](evidence/a2-left-cursor-20261112-01/user-game-test.go.txt)、[环境](evidence/a2-left-cursor-20261112-01/environment.log)、截图源与副本的相同 SHA256 以及源码运行前后哈希。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，GOWORK=off 只设于验证进程并还原；执行 `go test -mod=readonly ./internal/minesweeper -run '^TestLeftPlayerCannotMoveCursor$' -count=1 -v -timeout=30s`，1 PASS、退出 0，见[定向日志](evidence/a2-left-cursor-20261112-01/targeted-test.log)。源码前后哈希一致。
- 截图最后一次运行也显示同名用例通过；上方还残留一次 FAIL 尾部，但没有失败原因或当时源码，不能作为缺陷失败对照。本轮只运行该定向用例，未重跑全包/全库、构建、Race 或联机，A2 未完成。没有 AI 源码改动、提交、推送或共享文件修改。下一步以 TODO 骨架验证退出后的首次翻格被拒绝，再接续生命周期业务。

## A2 回到全员离开取消的业务实现（2026-11-13）

- 用户质疑是否要一直做防御测试，AI 确认此前相似练习拆得过细；用户随后要求继续。暂停尚未编写的退出后翻格测试，将下一步改为已授权 A2 的取消业务，不改变阶段范围或核心作者分工。
- 本轮重新读取共享规则及本地状态，静态追踪 loop.go 的 runOnce：GAME leave → Game.RemovePlayer → 记录离开日志，缺少 AllPlayersLeft 判断和退出本局；Loop 将 runOnce 返回 true 解释为关闭服务循环、false 解释为进入新大厅。gameOverCh 的现有消费者仅记录带玩家编号的胜者日志，不适合作为取消结果的编号占位，本步不使用它伪造胜负。
- 本轮仅提供 TODO 指导：Game 增加本局取消标记与 Cancel 方法，再在 GAME leave 分支接入取消和返回大厅。结束结果快照、固定名单及本步验证仍待实现，不将状态标记或日志等同持久保存结果。用户尚未填写，AI 未改业务/测试源码；本轮未运行程序测试、构建、联机、提交或推送。

## A2 用户取消草稿审查（2026-11-13，待修正）

- 用户按思路和 TODO 骨架保存 Canceled、Cancel 及 GAME leave 全员离开判断，AI 读取实际文件，未代写业务或核心断言。草稿的 Cancel 仅检查和设置 Canceled；leave 分支直接设置 Over，未调用 Cancel，随后 return true 会让外层 Loop 停止；日志 gameID 为 string，却使用 %d。记录为提示下尝试，尚未实现取消状态与返回大厅的验收行为。
- 批次 `a2-cancel-draft-20261113-01` 保留[用户 Game 草稿](evidence/a2-cancel-draft-20261113-01/user-game.go.txt)、[用户循环草稿](evidence/a2-cancel-draft-20261113-01/user-loop.go.txt)、[环境](evidence/a2-cancel-draft-20261113-01/environment.log)及源码运行前后 SHA256。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`；GOWORK=off 仅设于验证进程并还原。`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s` 退出 1，报告 `log.Printf format %d has arg gameID of wrong type string`，见[原始日志](evidence/a2-cancel-draft-20261113-01/package-tests.log)。没有用例进入 RUN，这属于 go test 默认 vet 阶段失败，不能作为取消业务的失败回归证据；源码前后哈希一致。
- AI 仅给出方法与分支定位、结束状态与取消原因的区别、返回值含义及格式类型提示，等待用户修正保存。未运行全库、构建、联机、Race，无源码修正、提交、推送或共享文件修改；原有通过记录不代表当前草稿通过。

## A2 取消草稿第二次核对（2026-11-13，仍待修正）

- 用户已将 Cancel 改为设置 Over 与 Canceled，GAME leave 调用该方法，日志改用 %v，可正常输出 string。仍有两处行为错误：Cancel 提前返回条件仍是 Canceled，已胜/负但未取消的游戏会被标为取消；全员离开仍返回 true，使外层 Loop 停止。AI 未修改源码，只提供具体状态和调用链提示供用户继续修正。
- 独立批次 `a2-cancel-review-20261113-02` 保留用户 Game / Loop 快照、环境及源码前后哈希。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`；验证进程 GOWORK=off 后还原。`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：12 个顶层 PASS、0 SKIP、0 FAIL、退出 0，见[包测试日志](evidence/a2-cancel-review-20261113-02/package-tests.log)及[计数](evidence/a2-cancel-review-20261113-02/test-summary.json)。源码运行前后 SHA256 一致。
- 本次证明日志格式检查已不再阻断、现有 12 项用例通过；用例没有调用新 Cancel 或 runOnce，不能证明取消和回大厅符合规则。未新增测试、未运行全库/构建/联机/Race，无提交、推送或共享文件修改，继续等待用户修正两处业务判断。

## A2 取消状态与重新打开大厅的有限诊断（2026-11-13）

- 用户第三次保存取消修改：Cancel 增加 Over 提前返回（原 Canceled 提前检查尚保留），GAME leave 返回 false。AI 未修改业务或测试源码。只读判断行为修正后，使用仓库外临时 Go 诊断源及 -overlay 执行真实 Game.Cancel / Loop；仅记录业务输出，不代写核心业务断言，测试框架只负责有界等待及清理。
- 批次 `a2-cancel-observe-20261113-01` 保留[用户 Game](evidence/a2-cancel-observe-20261113-01/user-game.go.txt)、[用户循环](evidence/a2-cancel-observe-20261113-01/user-loop.go.txt)、[AI 诊断源](evidence/a2-cancel-observe-20261113-01/diagnostic.go.txt)、[环境及临时路径](evidence/a2-cancel-observe-20261113-01/environment.log)及源码前后相同 SHA256。Go 1.27.3 / Windows amd64，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，GOWORK=off 仅设于验证进程并还原。
- 首次诊断因 AI 误用不存在的 protocol.MsgLobbyFull 常量而编译失败、退出 1，未执行观察；失败日志及源文件分别保留为 diagnostic-scaffold-error.log 和 diagnostic-before-correction.go.txt。这属于 AI 辅助代码错误，不是用户业务失败。核对协议后仅将诊断常量修正为 MsgLobbyState；重新执行 `go test -mod=readonly -overlay <仓库外 overlay.json> ./internal/minesweeper -run '^TestObserveCancelLifecycle$' -count=1 -v -timeout=10s`，退出 0，完整绝对路径与命令见[日志](evidence/a2-cancel-observe-20261113-01/diagnostic.log)。
- 实际观察：进行中的 Game 调用 Cancel 后 Over=true、Won=false、Canceled=true，重复调用不变；预置胜局 Over=true/Won=true 和败局 Over=true/Won=false 调用后 Canceled 均保持 false，重复调用不变。真实 Loop 按 Cooperative / Easy 两人开局，Alice 离开后 Bob 的光标请求仍产生广播；Bob 离开记录对局编号、restarting、lobby open，Carol 随后被接受并收到大厅广播；通过关闭 stop 清理循环。前两者为直接状态诊断，后者为本机通道驱动游戏循环，不含 WebSocket、渲染或手动联机。
- 诊断测试的 PASS 只表示观察程序正常执行，业务结论由上述输出人工核对，不冒充用户核心断言或正式失败→通过回归。本轮未重跑原有包/全库、构建、Race、性能或完整联机；未修改源码、提交、推送或共享文件。临时 overlay 留于环境日志所示系统临时目录，长期证据副本在本批次；不假定临时目录永久保留。下一步指导用户定义结束结果快照，保存及冻结结果仍未完成。

## A2 用户定义结果类型（2026-11-13）

- 用户按 TODO 骨架在 game.go 的 MinePlayer 后新增 MatchResult，独立填写 GameID string、Won bool、Cancel bool、PlayerResult []MinePlayer。AI 读取实际保存文件，确认所需数据及值列表类型齐全，沿用用户字段命名，未代写源码。
- 本步只有类型声明的静态核对，没有程序测试或结果生成行为，不把 []MinePlayer 的声明本身当作完成复制或冻结。下一步只提供 SnapshotResult 方法的 TODO 骨架和指针取值/独立切片提示，由用户实现复制；保存和跨局留存仍待接入。未提交、推送或修改共享文件。

## A2 用户实现快照复制方法（2026-11-13）

- 用户按思路与 TODO 骨架自行补齐 SnapshotResult：创建 res，复制 GameID/Won/Cancel；以 make([]MinePlayer, 0, len(g.Players)) 创建独立结果切片，逐项 append(*p)，赋给 PlayerResult 后返回。AI 读取实际文件，确认字段对应、切片长度/容量及值复制方式正确，未代写源码。当前 MinePlayer 仅含数值、bool 和 string，逐项值复制足以隔离其字段后续赋值，不泛化为任意嵌套结构的深复制。
- 批次 `a2-snapshot-method-20261113-01` 保留[用户源码](evidence/a2-snapshot-method-20261113-01/user-game.go.txt)、环境及源码前后相同 SHA256。Go 1.27.3 / Windows amd64，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，目录 `C:\Dev\Go\A`，GOWORK=off 仅设于验证进程后还原。`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：12 顶层 PASS、0 SKIP、0 FAIL、退出 0，见[日志](evidence/a2-snapshot-method-20261113-01/package-tests.log)。这些既有用例未调用新方法，因此仅为编译、默认静态检查及既有行为回归，不认定快照独立性已实测。
- 下一步仅安排一个与“结果固定”直接相关的快照隔离用例，由用户按 TODO 准备数据、取快照、修改原玩家并验证旧值保留；随后接入结束调用与跨局保存。未继续相似翻格练习，未执行全库/构建/联机/Race，无源码代写、提交、推送或共享文件修改。

## A2 用户验证快照玩家数据独立性（2026-11-13）

- 用户按 TODO 骨架自行补齐 TestSnapshotResultCopiesPlayers：创建 Alice，直接赋分 7，Cancel 后 SnapshotResult；再将原玩家分数改为 99 并离开，检查 result.PlayerResult[0] 的 Reveals 仍为 7、Left 仍为 false。AI 读取实际文件，未代写核心断言。分数断言的失败提示参数误用原游戏数值，已指出供用户修正；不影响当前条件本身判断快照值。
- 批次 `a2-snapshot-isolation-20261113-01` 保留[用户测试快照](evidence/a2-snapshot-isolation-20261113-01/user-game-test.go.txt)、环境及源码前后相同 SHA256。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，GOWORK=off 仅设于验证进程后还原。`go test -mod=readonly ./internal/minesweeper -run '^TestSnapshotResultCopiesPlayers$' -count=1 -v -timeout=30s`：1 PASS、退出 0，见[日志](evidence/a2-snapshot-isolation-20261113-01/targeted-test.log)。未注入错误取得失败对照，未重跑包/全库、构建、Race 或联机。
- 本次支持已生成快照与原玩家字段后续赋值隔离，不代表结束路径已调用、历史结果已留存或结果不可被持有者主动修改。下一步按 TODO 将 runOnce 的结束结果交回 Loop 并在进程内按局保存，不接数据库或磁盘功能。AI 未改源码、提交、推送或共享文件。

## A2 用户接入结束结果返回与保存，静态审查（2026-11-13，待调整顺序）

- 用户按 TODO 修改 runOnce 为 (bool, *MatchResult)，保留两个 stop 的 true/nil，胜负和取消分支返回 false 及快照地址；Loop 在 for 外声明 []MatchResult，接收两个返回值并在结果非 nil 时追加其值、记录编号和数量。AI 读取实际保存文件，确认这些部分及胜负分支等待前取快照的位置正确，未代写源码。
- 静态发现 GAME leave 全员离开分支先执行 SnapshotResult，随后才 Cancel，返回的快照 Cancel 标记因此仍为 false。已给用户“先完成结束状态转换、再取快照”的顺序提示，等待其自行调整。此前测试中的分数失败提示参数仍读取原游戏，亦提醒改为被检查的快照值。
- 本轮在运行前发现明确业务顺序错误，未运行新的测试/构建/诊断，不沿用旧测试宣称保存链路通过；尚无跨局保存的运行结论。未修改业务或测试源码，无提交、推送、共享文件或依赖修改。

## A2 保存路径观察及连续两局重号（2026-11-13）

- 用户修正为先 Cancel 再 SnapshotResult，并将快照测试的失败提示参数改为快照值；AI 读取实际保存文件，未改源码。独立批次 `a2-result-save-20261113-01` 保留用户 Loop / 测试快照、AI 临时诊断源、环境及源码前后相同 SHA256。
- Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`，验证进程 GOWORK=off 后还原。`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：13 PASS、0 SKIP、0 FAIL、退出 0，见[日志](evidence/a2-result-save-20261113-01/package-tests.log)及[计数](evidence/a2-result-save-20261113-01/test-summary.json)。
- 使用仓库外临时源及 overlay 执行 `go test -mod=readonly -overlay <临时 overlay.json> ./internal/minesweeper -run '^TestObserve(ReturnedCancel|TwoStoredMatches)$' -count=1 -v -timeout=10s`，完整绝对路径在[诊断日志](evidence/a2-result-save-20261113-01/diagnostic.log)和[环境](evidence/a2-result-save-20261113-01/environment.log)。AI 辅助诊断无核心业务断言，仅有输入/广播的有界等待与停止清理；退出 0 不等于所有观察值符合需求。
- runOnce 实际返回 stopped=false、Won=false、Cancel=true，快照有 Alice / Bob 两个玩家且 Left=true，支持取消结果顺序修正。真实 Loop 连续两局 Cooperative / Easy 全员离开后，保存日志 total=1、total=2，随后重新打开大厅接受 Erin，支持循环外列表累积。胜负快照分支本轮未执行；非零成绩保留由此前用户快照隔离用例另行支持。
- 同一诊断中两局保存编号均为 `mine-1794578159024-4825`，实际违反“再开局新编号”；不能据总计 2 或测试框架 PASS 认定 A2 验收通过。源码现有时间加随机数形式不能保证唯一，下一步在当前进程内结果寿命范围改用 Loop 递增编号并显式传入 runOnce / NewGame，保持布雷 rng 专用于游戏。仅给用户思路和 TODO 修改路径，尚未编写该修复；固定名单仍待处理。
- 本轮未运行全库、构建、Race、性能或 WebSocket/终端联机，无业务源码代写、提交、推送或共享文件修改。临时 overlay 的长期可查副本在本批次，系统临时目录不保证永久保留。

## A2 用户修复对局编号重复及本次提交前验证（2026-11-13）

- 用户要求提交推送后，又要求先核对并完成今天的编号修复。AI 重新读取实际文件，发现用户已经按 TODO 保存：Loop 的 matchNumber 在 for 外初始化、每轮递增，以 `mine - %d` 生成编号；通过新增 string 参数传给 runOnce、NewGame；Game 直接保存编号，不再调用时间及布雷 rng 生成标识；newTestMineGame 明确提供固定编号 `11`。核心业务与辅助构造器调用由用户填写，AI 未修改源码。
- 调用链为 Loop 生成编号 → runOnce 接收 → NewGame 保存至 MatchId → 结束分支 SnapshotResult 复制为 GameID → Loop 将结果值追加至进程内列表。选择单 Loop 顺序计数以满足当前单房间及进程内结果寿命，不增加依赖或全局并发计数；服务重启后编号可重复，不宣称跨进程全局唯一。
- 独立批次 `a2-id-fix-20261113-01` 保留用户三个源码快照、AI 诊断源、环境及源码前后相同 SHA256。Go 1.27.3 / Windows amd64，目录 `C:\Dev\Go\A`，验证前 HEAD `d02f62199e2b120a9044b3bdbd6abdc3e5beb25e`；GOWORK=off 仅设于验证进程后还原。复用前批次的相同诊断步骤，仅更新 runOnce 的新增编号实参以适配签名，没有给诊断添加核心业务断言。完整 overlay 命令在[诊断日志](evidence/a2-id-fix-20261113-01/diagnostic.log)。
- 实际观察：直接 runOnce 返回 stopped=false、指定 probe-match、Cancel=true、两位玩家 Left=true；真实 Loop 连续两局分别保存 `mine - 1` / total=1、`mine - 2` / total=2，第三个大厅接受 Erin。对照上一批次连续两局相同编号，支持本次进程内编号修复；诊断 PASS 只表示观察程序执行完成，编号结论来自原始输出核对，没有冒充用户编写唯一性断言。
- `go test -mod=readonly ./... -count=1 -v -timeout=60s`：42 顶层 PASS、0 SKIP、0 FAIL、退出 0，见[全库日志](evidence/a2-id-fix-20261113-01/full-tests.log)及[计数](evidence/a2-id-fix-20261113-01/test-summary.json)；`go build -mod=readonly ./...` 无输出、退出 0，见[构建日志](evidence/a2-id-fix-20261113-01/full-build.log)。未执行 Race、性能或本轮 WebSocket/终端联机，未新增胜负返回值的运行证据。
- 本次提交范围是自 d02f621 起的 A2 阶段进度、编号修复、项目约定及原始证据，不标记 A2 完成；固定名单及必要综合验收仍待后续。用户明确授权推送个人仓库 origin/main，已只读确认目标为 LearnedYet/multiplayer-terminal-games、远程 main 尚在原始基线 488dd313；将使用普通非强制推送，连同此前未推送的 A1 提交一并上传，不推送 upstream。实际提交/推送结果待 Git 命令回执核对，不预填成功；今天不再开始其他业务修改。
- 暂存核对：首次全量 `git diff --cached --check` 将原始 source-diff.patch 内必需的补丁上下文空格识别为空白问题而退出 2，提交命令尚未执行。保留这份原始证据不修改；后续空白检查仅排除 `docs/evidence/a2-handoff-check-20261112-01/source-diff.patch`，仍检查其余源码、文档及证据，不调整全局 Git 空白规则。

## A2 进度提交推送回执与换聊天交接（2026-11-13）

- 编号修复验证后，AI 按用户明确指令提交三个扫雷源码/测试文件、本地约定及原始证据，共 87 文件，提交 `5e47ebf6571456f1833e628e2bec136aeaf1aefc`，标题 `feat: preserve minesweeper match results and assign sequential IDs`。阶段未完成状态及贡献归属写入提交正文；没有改许可证、依赖或 BASELINE。
- `git push -u origin main` 退出 0，将个人仓库 LearnedYet/multiplayer-terminal-games 的 main 从 488dd31 快进到 5e47ebf，包含此前 d02f621；branch.main 改为跟踪 origin/main，未推送 upstream、未强推。随即 ls-remote 核对远程 main 与 HEAD 完全相同，工作区条目数 0。原回执存于 .git，换聊天时复制到[持久交接批次](evidence/a2-handoff-20261113-01/commit-and-push-receipt.log)。
- 用户准备换聊天后，本轮再次读取共享规则、本地状态、Git refs 及回执，三个扫雷源码/测试文件与已验证批次 SHA256 一致，见[交接检查](evidence/a2-handoff-20261113-01/handoff-check.log)。没有重跑测试、重新访问远程、开启新业务或再次提交推送。新聊天应从固定开局名单的实际连接/ID 接纳路径调查开始，按原定 TODO 分工继续 A2。
- 推送后的本次更新仅为 AGENTS、CURRENT_STATE、本文及独立交接回执批次，尚未提交；不能将文档未提交误判为有未验证核心代码，也不能沿用已经执行完的推送授权自动上传新任务。阶段综合验收、胜负结果完整链路联机及固定名单仍待后续，42 PASS 与构建结果沿用原始批次，未冒充本轮重测。

## A2 固定开局名单接续调查（2026-11-13）

- AI 只读核对项目及共享规则、状态、历史证据和真实源码。`git status --short` 显示 AGENTS、CURRENT_STATE、EVIDENCE 及 a2-handoff 回执目录；`git rev-parse HEAD origin/main` 两行均为 5e47ebf6571456f1833e628e2bec136aeaf1aefc；`git diff --exit-code HEAD -- cmd internal go.mod go.sum LICENSE` 无输出、退出 0。远程状态仅核对既有推送回执与本地 ref，本轮未联网、未执行测试或构建，42 PASS 仍归属此前编号修复批次。
- 静态调用链：cmd/mineserver/main.go 的 mineWSHandler → Hub.Register → Hub.Run 分配编号并发送 JoinMsg → runOnce 的 GAME join；输入由独立 readPump 直接进入 mineCh。Hub 在加入回复 false 时只发送错误消息，不阻止后续输入；注销释放编号后 assignID 可复用，三个事件通道没有跨通道先后保证。存在旧玩家 leave 尚未消费、新连接输入以相同编号到达的可能路径，尚未运行复现，不能仅用 Left 或拒绝 join 回复证明安全。
- 本轮只完成调查及本地文档记录，未改业务或测试源码。下一小步由用户填写 GAME join 拒绝逻辑，随后审查；连接接纳与输入放行仍需后续处理，A2 未完成，无提交、推送、共享文件修改。

## A2 用户拒绝 GAME 中途加入并完成包回归（2026-11-13）

- 用户按 TODO 保存 internal/minesweeper/loop.go 的 runOnce GAME join 分支：删除按模式中途 AddPlayer、接受回复和棋盘广播，只回复 false；大厅 join 与其余业务未改。AI 读取实际 diff，确认未返回退出本局，未代写业务或核心断言。
- 批次 [a2-fixed-roster-join-20261113-01](evidence/a2-fixed-roster-join-20261113-01/) 保存用户 loop 快照、环境、测试原始输出及三个扫雷源码/测试前后 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，验证进程 GOWORK=off 后还原；执行 `go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：13 PASS、0 SKIP、0 FAIL、退出 0，源码哈希一致。现有用例未调用 runOnce，不将包通过记为中途加入或连接身份验收；本轮未执行全库、构建、Race 或联机。
- 本步学习点是通过回复通道拒绝接纳，同时保留正在运行的对局。连接仍可发送输入、编号可复用的问题尚未解决；下一步先讨论服务端连接身份与拒绝即断开的取舍，重大方案确认后继续 TODO 带学。无提交、推送、依赖/许可证或共享文件修改。

## A2 连接身份方案批准与第一步准备（2026-11-13）

- 用户明确同意“增加服务端连接身份”的推荐方案，保留已有连接等待/再加入行为；不选择拒绝加入即断开的备选。核心业务及核心断言继续由用户按 TODO 小步实现，AI 调查、讲解、审查及验证，不将方案批准当作 AI 代写授权。
- AI 静态核对 Hub.NewClient 由 HTTP 处理函数调用，后续才 Register 并启动读循环。第一步拟在 NewClient 固定身份，Hub 使用原子计数器支持并发创建连接；用本机 `go doc sync/atomic.Uint64`、`go doc sync/atomic.Uint64.Add` 核对零值及返回值含义，均正常退出。没有修改业务/测试源码或执行新测试、构建、联机；上一批 13 PASS 不证明尚未实现的连接身份方案。
- 当前仅提供 hub.go 的字段、分配及读取方法 TODO，等待用户保存；后续沿相关内部消息传递身份，并在大厅/本局操作与离开路径匹配开局连接。客户端线协议、账号系统和跨进程持久性不在方案内，无提交或推送。

## A2 用户保存连接身份分配与 server 包回归（2026-11-13）

- 用户按 TODO 保存 internal/server/hub.go：Client 的 sessionID 为 uint64，Hub 的 sessionCounter 为 atomic.Uint64，NewClient 将 Add(1) 返回值保存到新连接，SessionID() 只返回已保存值。AI 读取实际 diff 与源码，确认计数器在 Hub 存续、分配在新连接创建时完成、getter 不再递增；未修改业务或核心断言。
- 批次 [a2-session-allocation-20261113-01](evidence/a2-session-allocation-20261113-01/) 保留用户 hub 快照、环境、原始输出及 hub.go/hub_test.go 前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，验证进程 GOWORK=off 后还原；执行 `go test -mod=readonly ./internal/server -count=1 -v -timeout=30s`：5 PASS、0 SKIP、0 FAIL、退出 0。仅支持编译及既有 Hub 注册/编号复用/广播/容量/模式转发回归；测试未新增对 session 的断言，不认定已验证并发身份唯一性或完整连接隔离。未执行全库、构建、Race 或联机。
- 下一步用户填写相关内部消息的 SessionID 字段及来源，批量处理直接的重复赋值；Start 消息的函数参数传递随后处理，最后再接入本局名单匹配与必要业务验证。无提交、推送、依赖/许可证或共享文件修改。

## A2 消息身份传递草稿审查（2026-11-13，待补字段）

- AI 读取用户实际 diff：MineInputMsg、扫雷 SetModeMsg 已声明 SessionID uint64；三类操作和模式消息使用 c.SessionID()，Hub register/rejoin 两处 JoinMsg 和 unregister 的 LeaveMsg 均使用 c.sessionID，未重新分配或读取客户端载荷。internal/snake/loop.go 实际保存的 JoinMsg、LeaveMsg 尚无 SessionID 字段，是当前明确缺项；AI 未代写字段或继续下一个依赖步骤。
- 独立批次 [a2-session-messages-draft-20261113-01](evidence/a2-session-messages-draft-20261113-01/) 保存 Hub 与 Snake 消息定义快照、环境及四个相关源码前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅设于验证进程并还原；`go test -mod=readonly ./internal/server -count=1 -v -timeout=30s` 退出 1：hub.go 的两处 JoinMsg、单处 LeaveMsg 报 unknown field SessionID。编译失败，无用例进入 RUN，不记测试用例断言失败或通过。
- 当前仅提示用户在 internal/snake/loop.go 的 JoinMsg、LeaveMsg 中补齐或保存字段，已有赋值不改。未执行其他包/全库测试、构建、Race 或联机；无源码代写、提交、推送、依赖/许可证或共享文件修改。

## A2 用户补齐消息身份字段与相关四包回归（2026-11-13）

- 用户保存 internal/snake/loop.go 中 JoinMsg、LeaveMsg 的 SessionID uint64，AI 读取实际文件确认此前缺项已补齐。该文件伴有整文件格式化；通过移除两个新字段后，将当前文件和 HEAD 分别送入 gofmt 并比较，结果相同，见本批次 format-check.log。未改 Snake 玩法逻辑，AI 未代写业务/核心断言。
- 独立批次 [a2-session-messages-20261113-01](evidence/a2-session-messages-20261113-01/) 保存用户 Snake loop 快照、环境、原始测试输出和五个相关源码/测试前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅设于验证进程后还原；`go test -mod=readonly ./internal/server ./internal/snake ./internal/minesweeper ./cmd/mineserver -count=1 -v -timeout=30s`：30 顶层 PASS、0 SKIP、0 FAIL、退出 0。上一草稿的三处 unknown field 编译失败已消除，既有测试无新增断言。
- 已有 TestMineCursorUsesConnectionIdentity 使用真实本机 WebSocket，输出收到 PlayerID=0、SessionID=1、Action=2、X=3、Y=4，支持一个正常操作消息携带连接身份的有限观察；其核心断言仍只验证既有玩家编号与坐标，不冒充 session 唯一性或复用隔离回归。未执行全库、单独构建、Race、手动终端联机或性能测试。
- 下一步指导 StartMsg → ForwardStart 参数 → 三处调用的身份传递，属于既定方案；大厅和 GAME 尚不检查 session，固定名单整体未完成。无提交、推送、依赖/许可证或共享文件修改。

## A2 用户接通开局身份参数与四包回归（2026-11-13）

- AI 读取实际保存内容：StartMsg 新增 SessionID uint64；ForwardStart 接收 sessionID 参数并填入 StartMsg；mineWSHandler 的 MsgStartGame、共享 Client.readPump 的 MsgStartGame/MsgNextLevel 三处调用均传递当前连接的已有身份。client.go 同时格式化，未发现其他行为变更；AI 未修改源码或核心断言。
- 批次 [a2-session-start-20261113-01](evidence/a2-session-start-20261113-01/) 保存四个用户文件快照、环境、原始输出及五个相关源码前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅设于验证进程后还原；执行 `go test -mod=readonly ./internal/server ./internal/snake ./internal/minesweeper ./cmd/mineserver -count=1 -v -timeout=30s`：30 顶层 PASS、0 SKIP、0 FAIL、退出 0。验证调用签名及既有行为回归，尚无新增 StartMsg.SessionID 核心断言或固定名单身份匹配，不把通过升级为阶段验收完成。
- 下一步由用户在 runOnce 记录大厅接纳的连接身份，匹配大厅的离开/模式/开局消息；GAME 的身份判断仍待后续。未执行全库、单独构建、Race、手动终端联机或性能测试，无提交、推送、依赖/许可证或共享文件修改。

## A2 大厅身份匹配草稿审查（2026-11-13，待调整广播顺序）

- 用户按 TODO 在 runOnce 创建 playerSessions map，成功加入后写入身份并刷新名字；大厅 leave、setMode、start 先检查编号存在且身份相同，否则 continue；有效 leave 清理 lobbyPlayers 与 playerSessions。AI 读取实际源码和 diff，确认这些判断及原房主限制正确，未代写或修改业务。
- 静态发现大厅 join 的 if ok 中先 broadcastLobby，随后才更新两个 map。已有编号被新连接复用且旧 leave 尚待处理时，addToLobby 对已有编号返回 true，广播会读到旧名字，随后更新无即时广播。当前仅提示用户将两个赋值移到广播前；未运行新测试或将上一批 30 PASS 冒充草稿验收。
- 等待用户保存顺序修正后验证，再接 GAME 身份匹配；无提交、推送、依赖/许可证或共享文件修改。

## A2 用户修正大厅广播顺序与扫雷包回归（2026-11-13）

- AI 读取实际保存的 runOnce：LOBBY join 成功后先更新 playerSessions 与 lobbyPlayers，再 broadcastLobby，原先广播旧名字的顺序问题已修正；大厅三处身份匹配及有效 leave 的两表清理仍在。AI 未改业务或核心断言。
- 批次 [a2-session-lobby-20261113-01](evidence/a2-session-lobby-20261113-01/) 保存当前 loop 快照（尚无 GAME 身份匹配，可供后续相同回归做旧实现对照）、环境、原始输出及三个扫雷源码/测试前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅设于验证进程后还原；`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：13 PASS、0 SKIP、0 FAIL、退出 0。既有测试未调用 runOnce，大厅身份规则仅静态审查，运行验证待必要业务回归，不把该结果记为阶段验收。
- 下一步用户按 TODO 在 GAME 的输入与离开分支匹配开局身份；名单身份表在 GAME 中保持固定，Game.Left 继续负责停止离开者操作。未执行全库、单独构建、Race 或联机，无提交、推送、依赖/许可证或共享文件修改。

## A2 用户接入 GAME 身份匹配及综合回归外壳准备（2026-11-13）

- AI 读取用户保存的 GAME mineCh/leaveCh：先从 playerSessions 查询身份并匹配，不存在或不同则 continue，匹配后才执行操作或 RemovePlayer；GAME 中不修改该表，既有取消和快照处理保留。业务由用户按 TODO 编写，AI 未代改。
- 批次 [a2-session-game-20261113-01](evidence/a2-session-game-20261113-01/) 保存用户 loop 快照、环境、日志及三个既有扫雷源码/测试前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅设于验证进程后还原；新增外壳前执行 `go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：13 PASS、0 SKIP、0 FAIL、退出 0。该批既有用例不调用 runOnce，不能替代新场景回归。
- AI 准备 internal/minesweeper/loop_test.go：TestRunOnceRejectsReusedConnection 顶部仅有观察变量、日志、业务 TODO 与故意失败的占位；exerciseReusedConnection 负责通道、两人 Cooperative / Easy 开局、复用 ID 的 join/Flag/leave 输入、顺序边界及清理。操作与离开顺序由 mineCh FIFO、无缓冲 leaveCh 接收确定，不靠短 sleep。此为直接驱动真实循环的内部场景，未使用 Hub/WebSocket，不将其记为真实网络复用测试。
- 外壳 `go test -mod=readonly ./internal/minesweeper -run '^$' -count=1 -timeout=30s` 编译成功、退出 0、无用例运行。随后 `go test -mod=readonly ./internal/minesweeper -run '^TestRunOnceRejectsReusedConnection$' -count=1 -v -timeout=15s` 到达占位 t.Fatal、退出 1；输出 lateAccepted=false、endedEarly=false，Cancel=true，结果保留 Alice/Bob 两人，Alice.Flags=1、Bob.Flags=0、两人 Left=true。仅支持外壳可运行的有限观察，未填核心断言的 FAIL 不是业务缺陷回归证据，不记新用例通过。
- 用户下一步填核心判断后，再用相同测试对比 a2-session-lobby 批次的未加 GAME guard 快照。未执行全库、单独构建、Race、手动联机或性能验证，无提交、推送、依赖/许可证或共享文件修改。

## A2 用户完成编号复用综合断言与失败通过对照（2026-11-14）

- AI 读取用户保存的 TestRunOnceRejectsReusedConnection：拒绝 lateAccepted、拒绝 endedEarly、检查 result 非 nil/Cancel/人数 2，再按 Name 识别 Alice/Bob、检查 Alice.Flags=1，拒绝其他名字并确认两人均找到。占位失败已删除；核心断言由用户按 TODO 填写，AI 没有代写，外壳仍归 AI 辅助。不认定用户独立设计整套并发测试。
- 独立批次 [a2-session-regression-20261114-01](evidence/a2-session-regression-20261114-01/) 保存用户测试快照、环境、overlay 映射、两次原始输出及五个相关源码/测试前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 仅在验证进程后还原。旧实现通过仓库外临时 overlay 使用 a2-session-lobby 批次保存的 user-loop-before-game-guard.go.txt，不修改工作区源码，临时路径和精确命令见 environment.log 与 before-targeted.log。
- 同一测试命令主体 `go test -mod=readonly ./internal/minesweeper -run '^TestRunOnceRejectsReusedConnection$' -count=1 -v -timeout=15s`，旧实现额外传入 `-overlay <本批次记录的路径>`：旧快照 lateAccepted=false、endedEarly=true，在用户的提前结束断言失败，退出 1；当前实现 lateAccepted=false、endedEarly=false，Cancel=true，名单 Alice/Bob，Alice.Flags=1，1 PASS、退出 0。旧快照缺少 GAME 输入和离开身份检查，其拒绝 join 及大厅匹配已存在，不把对照结论扩大为所有固定名单路径均获失败覆盖。
- 结论为真实 runOnce 内部通道场景的回归通过，未经过 Hub/WebSocket，不等同网络编号复用与慢连接离开的综合验收。本轮未执行全库、单独构建、Race、手动联机或性能验证。下一步用户收尾 Hub 的 stale rejoin 指针匹配及慢连接 LeaveMsg，以落实已有固定名单/离开规则；无提交、推送、依赖/许可证或共享文件修改。

## A2 Hub 收尾草稿静态审查（2026-11-14，待整理原检查与通知位置）

- AI 读取用户保存的 Hub.Run：rejoin 新增 existing/ok 查询及指针相等判断，符合既定方案；但旧的仅 exists 检查仍紧随其后，当前为重复查询，提示用户删除原检查块。慢连接 LeaveMsg 已携带 c.id/c.sessionID，字段正确；发送目前在关闭通道的匿名函数内部，提示按原要求移动到 }() 调用后、慢连接 if 内，关闭函数仅负责原有关闭动作。
- 此轮为源码及 diff 静态审查，AI 未代写核心业务；尚未运行 Go 测试、构建或连接诊断，不把已有回归结果记为这版 Hub 的验证。待用户保存两处局部整理后继续综合验收，无提交、推送、依赖/许可证或共享文件修改。

## A2 Hub 收尾与必要综合验收（2026-11-14，学习复盘待完成）

- 用户保存 Hub.Run：移除 rejoin 旧重复查询，保留存在性及 Client 指针检查；慢连接删除和关闭后发送带 c.id/c.sessionID 的 LeaveMsg。AI 读取实际文件，未代改业务。关于通知在匿名函数内/外，AI 已更正其为正常路径等效的组织取舍，不应将移动本身宣称为功能缺陷修复。
- 批次 [a2-acceptance-20261114-01](evidence/a2-acceptance-20261114-01/) 保存用户 Hub、诊断代码、临时工厂 hook 完整副本、overlay 映射、环境、原始输出及 cmd/internal 全部 Go 文件前后相同 SHA256。Windows amd64 / Go 1.27.3，HEAD 5e47ebf，GOWORK=off 只设于验证进程后还原；辅助源码、overlay 均在 environment.log 所列仓库外 Temp 路径生成，长期副本在本批次，生产源码未被 overlay 替换写入。
- Hub 有限观察：真实 Hub 把编号 0 从 alice/session1 复用给 carol/session2；同一 rejoin 队列依次发送旧/当前 Client 后，第一个转发的 join 为 carol/session2。向未消费发送队列的连接连续广播至饱和，收到其 PlayerID=0/SessionID=2 的 leave；后续 stale Unregister 被忽略，dave/session3 仍能复用编号加入和离开。见 diagnostic.log；这是辅助诊断输出核对，非用户新增核心断言。
- 本机网络观察：httptest/真实 WebSocket → mineWSHandler → Hub/Loop，Alice/Bob Cooperative / Easy 开局；Alice 断开，Carol 获得同号 0、SessionID=3 且收到游戏进行中错误。Carol 对 (0,0) 插旗后，Bob 对 (1,0) 插旗，观察 flagsLeft=9、前者隐藏/后者插旗；Bob 断开后第一局 Cancel=true，名单 Alice/Bob，Bob.Flags=1、两人 Left=true。Carol 同一连接下一局 PlayAgain，完成胜局，再次 PlayAgain 后踩雷完成败局；之后 Dave 进入大厅。
- 胜负诊断方法及限制：只在临时 Game 工厂中把新 Game 指针发送到诊断通道，原 RNG、布雷和玩法未修改；等待广播后读取棋盘选择未翻开的安全格或雷格，每次输入等待状态回复，不依赖猜测随机雷区。三次观察的快照分别为 mine - 1（取消，Alice/Bob）、mine - 2（胜局，Carol.Reveals=71）、mine - 3（败局，Carol.Reveals=46）；后局结束时再次观察的前局数据保持原值。观察快照来自 Game.SnapshotResult，真实 Loop 的生产日志依次保存编号 mine - 1/2/3、total=1/2/3；没有直接读取未暴露的进程内结果列表，快照隔离的正式断言由此前用户用例支持。
- 精确诊断命令见 diagnostic.log：`go test -mod=readonly -overlay <临时 overlay.json> ./internal/server ./cmd/mineserver -run '^TestObserveA2' -count=1 -v -timeout=45s`，两个辅助观察正常结束、退出 0；不将诊断框架 PASS 冒充用户核心断言或长期自动回归。使用真实本机网络但无手动终端画面，不属于公共部署。
- 正式全库验证：`go test -mod=readonly ./... -count=1 -v -timeout=60s` 为 43 顶层 PASS、0 SKIP、0 FAIL、退出 0，见 full-tests.log/test-summary.json；`go build -mod=readonly ./...` 无输出、退出 0，见 full-build.log；源码哈希全一致，见 hash-check.log。未执行 Race 或性能测试，无并发安全/吞吐性能结论。
- 本轮完成工程验证，仍待一次用户学习复盘；上游架构、用户在 Gemini 下早期实现、当前 TODO 业务与断言、AI 诊断外壳的归属分别保留。不新增相似练习，不提交、推送或开启新阶段；未改依赖、许可证、共享文件或其他项目。
- 哈希核对说明：首次直接比较 JSON 文本因 rg 枚举文件顺序不同而得到 false；随后按 Path 排序并比较 Path/Hash，差异 0，确认所有源码字节未变。原输出及修正核对均保留在 hash-check.log，不据枚举顺序差异推断源码变化。

## A2 学习复盘反馈（2026-11-14，接纳与连接存续待澄清）

- 用户回答连接身份作用为“防止新加入的玩家操作已经掉线的玩家”，回答重启结果为“目前没有数据持久化，所以重启就没了”；方向及保存边界正确，另需补充编号计数也会重置。用户对拒绝加入后仍检查操作/leave 明确表示未明白，不记为已掌握或阶段整体完成。
- AI 补讲已选方案保留连接等待下一局：GAME join 返回 false 只阻止 AddPlayer，Hub 发送错误提示但不停止 readPump；操作仍可进入 mineCh，断开仍会 Unregister 并发送 leave。用同号不同身份的 Carol leave 错标 Alice、导致旧回归提前取消的实际场景对应解释。无源码或测试修改，无新测试/构建、提交推送或新阶段；等待该概念的简短说明，不追加相似练习。

## A2 学习收尾与阶段交付（2026-11-14，完成停止）

- 用户复盘说明连接身份用于防止新连接影响原玩家记录、未持久化导致重启丢失；对拒绝 join 与连接存续的区别，经实际失败回归场景讲解后，回答“因为拒绝加入复用了编号，也被监测，会产生消息，所以要加身份识别”。核心理解为拒绝本局接纳后连接仍产生操作/离开事件，必须识别连接；AI 同时澄清 Hub 先分配复用编号，之后游戏循环拒绝接纳的顺序。该概念在指导和解释下完成说明，不认定完全独立设计模块或熟练掌握所有并发行为，不追加练习。
- A2 范围内的固定开局名单、离开停止操作并保留成绩、全员离开取消、胜负/取消快照及进程内保存、再次开局递增编号完成；必要工程验证及学习复盘已完成，阶段于 2026-11-14 完成停止。结果仅进程内保留、重启丢失，编号仅本次 Loop 内唯一；诊断辅助/正式断言、网络自动观察/手动体验等边界沿用上节，不新增未测性能或 Race 结论。
- 当前接续的源码/测试文件为 cmd/mineserver/main.go、internal/minesweeper/loop.go、internal/minesweeper/loop_test.go、internal/server/hub.go、internal/server/client.go、internal/snake/loop.go；更早的 game.go/game_test.go 生命周期成果已包含在 5e47ebf 中。调用链：NewClient 分配身份 → 内部事件携带身份 → runOnce 大厅接纳并记录 → GAME 固定身份校验 → Game 离开/结束及 SnapshotResult → Loop 追加结果；复用原架构和线协议，不引入账号、持久化或通用抽象。
- 贡献归属：上游提供 Hub、游戏循环、协议、玩法与客户端；早期生命周期代码由用户在 Gemini 指导下编写；当前用户按 TODO 完成核心业务和代表性断言，AI 提供讲解、调查、测试外壳、临时诊断、验证及本地文档。不得将全部上游能力或 AI 外壳记为用户独立开发。
- 收尾只读核对：21 个 cmd/internal Go 文件与 a2-acceptance 批次记录的 Path/Hash 差异为 0；HEAD 和本地 origin/main 仍为 5e47ebf6571456f1833e628e2bec136aeaf1aefc。未重跑测试/构建、未重新联网核对远程；43 PASS 和构建通过引用该验收批次。新源码/文档/证据仍在工作区，无提交、推送、部署或新阶段，不修改 BASELINE、依赖、许可证或共享文件。

## A2 完成后的提交推送准备（2026-11-14）

- 用户明确要求“提交推送github”，授权本次 A2 接续源码、文档及证据提交到个人 origin/main。AI 核对当前 main 分支、origin 为 LearnedYet/multiplayer-terminal-games，upstream 为 siyad01 上游，本次仅操作 origin。
- 提交前源码与 a2-acceptance 批次按 Path/Hash 比较差异为 0，沿用已完成的 43 PASS 与构建通过，不重复运行测试或修改业务。LICENSE、go.mod、go.sum、BASELINE 相对 HEAD 无差异。提交范围包含此前未提交的交接材料、六个源码/测试文件、本地阶段文档及本轮独立证据批次；提交/推送结果以实际 Git 回执为准，不预填成功。

## A2 最终代码提交推送回执（2026-11-14）

- 按用户“提交推送github”指令，提交 75 个文件至 e3413e3b66deaac9113f4ab1e94cd77dcbe59b18，标题 `feat: freeze minesweeper match rosters by connection identity`；包括六个接续源码/测试文件、阶段完成文档、交接材料和原始验证批次。提交正文保留用户/Gemini/AI 各自贡献和验证边界；没有修改 LICENSE、依赖、BASELINE 或共享文件。
- 原始暂存空白检查在旧用户源码快照 docs/evidence/a2-session-messages-draft-20261113-01/user-snake-loop.go.txt 发现一处原始行尾空格。保留证据内容，仅将该文件排除后检查其余暂存内容通过，不修改 Git 全局空白规则。源码与验收哈希一致，沿用 43 PASS 与构建结果，本轮未重跑测试。
- GitHub 直连 ls-remote 低速超时，旧代理端口 7897 已无监听；本轮核对系统当前代理及 netstat，发现 127.0.0.1:10808 在监听。仅对 Git 命令指定该代理并保留 http.sslVerify=true，ls-remote 与 push 均成功，未修改持久代理配置或跳过现有 hooks/LFS。
- `git push origin main` 普通非强制推送退出 0，个人 origin/main 从 5e47ebf 快进至 e3413e3。随后联网 ls-remote 与本地 HEAD 完全相同，回执归档前工作区条目数为 0；原始输出与核对见 [本批次](evidence/a2-push-20261114-01/)。未推送 upstream、未部署或开启新阶段。
- 本段、项目约定和状态更新及原始回执随后通过独立文档提交归档并推送；最新管理提交由 Git refs 和最终操作回执确认，避免在提交内自引用自身哈希。本次授权完成后不延续为以后自动提交推送。

## 单项贡献模板

复制后填写；未发生的字段留空。

```text
记录编号：
日期 / 执行者：
任务范围与验收条件：
上游来源 / 复用部分：
基准提交 / 完成提交或补丁：
改动文件：
用户独立实现与理解的部分：
工具 / AI 辅助的部分：
关键调用链（入口 → 处理 → 状态 / 输出）：
设计理由 / 已批准的方案：
实际测试命令 / 环境 / 输入：
实际结果 / 退出码 / 失败原因：
原始日志或附件位置：
未验证事项 / 局限：
需要学习的知识 / 用户复述或独立修改证据：
```

## 量化测试模板

对比必须使用相同机器、工具链、构建参数和工作负载；明确指标定义、单位、预热、重复次数与统计方法。记录原始数据和失败样本，不把一次冒烟、游戏得分或编译用时当作服务端性能；不把跳过测试算作通过。

| 实测日期 | 指标与单位 | 工作负载 / 样本量 | 环境与命令 | 基准提交 / 原始值 | 完成提交 / 原始值 | 统计方法 / 结果 | 原始证据 / 限制 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| | | | | | | | |

## 个人贡献索引

以下只记录实际发生的贡献与验证，并标明指导与独立修改的范围；Phase 0 文档成果仍按上面的事实说明归属。

| 编号 | 任务 | 基准 / 完成提交 | 用户贡献 | 上游复用 / 辅助 | 证据位置 |
| --- | --- | --- | --- | --- | --- |
| A1.1-R1（指导下完成） | 字段赋值及消息断言，自动与有限联机验证通过 | `488dd313bc13` / `d02f621` | 补字段赋值、玩家 1 输入及两字段判断；格式化报告经指导修正 | 上游架构复用；AI 讲解、格式整理及验证辅助 | [回归通过](evidence/a11-regression-20261111-01/current-targeted-test.log)、[联机截图](evidence/a11-snake-after-fix-20261111-01/server-mode-2-and-clean-shutdown.png) |
| A1.2-R4（指导下完成，已交付） | 光标身份来源修复，自动验证及正常位置同步有限观察完成 | `488dd313bc13` / `d02f621` | 按示例保存身份断言及 `c.ID()` 单行修复；本人操作双客户端并正确解释身份与坐标来源 | 上游协议/游戏复用；AI 提取处理函数、搭建网络外壳、讲解及验证 | [修复前失败](evidence/a12-cursor-before-fix-20261111-01/targeted-test.log)、[修复后通过](evidence/a12-cursor-after-fix-20261111-01/targeted-test.log)、[联机画面](evidence/a12-cursor-manual-20261111-01/02-two-clients-after-bob-move.png) |
| A2（指导下完成，阶段已结束） | 单对局生命周期、固定名单及连接身份隔离 | `d02f621` / `5e47ebf`、`e3413e3`，已推送个人 origin/main | Gemini 下早期实现；本聊天 TODO 核心业务、编号复用断言及提示下概念复盘 | 上游架构复用；AI 讲解、测试外壳、诊断及验证 | [复用回归对照](evidence/a2-session-regression-20261114-01/)、[综合验收](evidence/a2-acceptance-20261114-01/)、[推送回执](evidence/a2-push-20261114-01/) |
