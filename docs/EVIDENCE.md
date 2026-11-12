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
| A1.1-R1（指导下完成） | 字段赋值及消息断言，自动与有限联机验证通过 | `488dd313bc13` / 未提交 | 补字段赋值、玩家 1 输入及两字段判断；格式化报告经指导修正 | 上游架构复用；AI 讲解、格式整理及验证辅助 | [回归通过](evidence/a11-regression-20261111-01/current-targeted-test.log)、[联机截图](evidence/a11-snake-after-fix-20261111-01/server-mode-2-and-clean-shutdown.png) |
| A1.2-R4（指导下完成，已交付） | 光标身份来源修复，自动验证及正常位置同步有限观察完成 | `488dd313bc13` / 未提交 | 按示例保存身份断言及 `c.ID()` 单行修复；本人操作双客户端并正确解释身份与坐标来源 | 上游协议/游戏复用；AI 提取处理函数、搭建网络外壳、讲解及验证 | [修复前失败](evidence/a12-cursor-before-fix-20261111-01/targeted-test.log)、[修复后通过](evidence/a12-cursor-after-fix-20261111-01/targeted-test.log)、[联机画面](evidence/a12-cursor-manual-20261111-01/02-two-clients-after-bob-move.png) |
