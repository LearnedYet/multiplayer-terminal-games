# 当前状态

## 阶段与边界

- A1.2 扫雷光标消息身份归属修复于 2026-11-11（America/Los_Angeles）由用户开启并于当日完成停止。用户指导下保存核心断言及 `c.ID()` 修复，同一测试修复前失败、修复后通过，全库 38 项测试及构建通过，双客户端正常位置同步和停机有有限截图证据；简短学习复盘已完成。下一任务尚未授权。
- A1.1 模式切换修复与回归测试于 2026-11-11 完成交付；用户于 2026-11-02 开启任务，原任务书为 `C:\Dev\Go\docs\A_PHASE1_TASK.md`。其成果及证据保留，不重做验收。
- Phase 0 原始环境与基线于 2026-10-29 记录，本人最小复验于 2026-11-02 收尾；历史结果不重写。
- 原始上游与当前 HEAD：`488dd313bc13de387d12a24f5d05a5c28bed5899`；本地修改尚未提交。原始环境、结果及模块说明见 [BASELINE.md](BASELINE.md)。
- Snake、Minesweeper、WebSocket Hub、二进制协议及终端客户端均为上游已有能力。个人改动、AI 辅助及原始验证分别记录于 [EVIDENCE.md](EVIDENCE.md)。
- A1.1 未修改依赖、许可证或共享文件，未提交、推送或部署。A1.2 同样没有此类操作授权；构建产物与临时副本仍放仓库外。
- 用户说明 PORTFOLIO.md 为旧 Gemini 同步材料，已不再作为当前进度依据或同步待办。当前进度以本项目状态和证据为准，共同学习规则继续遵守。

## A1.2 交付（指导下完成）

- 范围与问题：仅处理扫雷光标操作的身份归属。修复前 `cmd/mineserver/main.go` 的 `MsgMineCursor` 使用消息中的 `data[1]`，真实 WebSocket 定向测试复现 Alice=0 的连接发送载荷 Bob=1，入口实际转发 PlayerID=1。用户已将此行改为 `c.ID()`，同一测试现在实际转发 PlayerID=0，动作与坐标仍正确。
- 目标：光标操作以服务端分配给发送连接的玩家编号为准，客户端填写其他玩家编号也不能移动别人的光标；正常客户端的光标同步继续可用。保留现有消息格式，不扩大协议或账号系统。
- 核心学习点：连接所代表的玩家与消息内声称的玩家之间的区别；外部输入的信任边界；用回归测试观察真正的入口行为。
- 核心作者及 AI 范围：核心修复和代表性断言由用户在指导下小步编写并保存；AI 负责只读调查、讲解、必要测试辅助、执行验证和本地文档。不是授权 AI 代写核心逻辑或整份测试答案。
- 调用链：客户端 `runMine` 中调用 `protocol.EncodeCursor` → WebSocket → `cmd/mineserver/main.go` 的 `MsgMineCursor` 分支 → `mineCh` → `internal/minesweeper/loop.go` 的 `ActionCursor` → `Game.MoveCursor` → 广播光标位置。
- 验收终点：同一回归用例在修复前失败、修复后通过，覆盖连接身份与消息编号不一致的情况，并确认坐标正常传递；相关测试、全库测试及构建按变更执行，双客户端正常光标同步有有限冒烟证据。未执行不记通过，不新增性能或完整规则结论。
- 最小测试入口：AI 先将 `cmd/mineserver/main.go` 原有 `/ws` 匿名处理函数原样提取为 `mineWSHandler`，生产和测试使用同一函数，提取时保留 `PlayerID: data[1]`；函数体在格式化前逐字一致。之后的核心修复由用户保存，修复前后归一化比较只有 `data[1]` → `c.ID()` 一行变化，测试字节哈希相同。
- 测试外壳：AI 新增 `cmd/mineserver/main_test.go`，使用本机临时 HTTP/WebSocket 服务和真实 Hub，先从 Welcome 获取 Alice、Bob 的实际编号，由 Alice 发送填有 Bob 编号的光标消息，再从同一 `mineCh` 接收 `msg`。连接、收发、动作与坐标检查为 AI 辅助；用户已按提供的三行示例，将 TODO 替换为 `msg.PlayerID != aliceID` 及 Fatalf 报告。没有启动完整扫雷循环或交互服务。
- 实际学习结果：用户对伪造编号会移动 Bob 光标的局部预测正确，按明确示例保存核心断言和修复，并操作双客户端回传截图。收尾时能用自己的话解释服务器分配玩家编号、信任消息编号会允许冒用身份，而坐标是客户端请求的位置；AI 补充坐标仍属外部输入，须按规则校验，本任务未新增坐标校验。记录为指导下完成及局部解释正确，未验证独立设计网络测试、独立修复或完整安全治理能力，不将聊天复述升级为独立掌握。
- 自动验证：同一 `go test -mod=readonly ./cmd/mineserver -run '^TestMineCursorUsesConnectionIdentity$' -count=1 -v -timeout=20s` 修复前因期望 0、实际 1 失败（退出 1），修复后收到 `{PlayerID:0 Action:2 X:3 Y:4}` 并通过（退出 0）。全库 `go test -mod=readonly ./... -count=1 -v -timeout=60s`：38 个顶层 PASS、0 SKIP、0 FAIL，5 个有测试包、3 个无测试包，退出 0；`go build -mod=readonly ./...` 无输出、退出 0。见 [修复前日志](evidence/a12-cursor-before-fix-20261111-01/targeted-test.log)、[修复后日志](evidence/a12-cursor-after-fix-20261111-01/targeted-test.log)及 [全库日志](evidence/a12-cursor-after-fix-20261111-01/full-tests.log)。未执行 Race，完整游戏规则与性能未验证。
- 本人联机：用户回传两张双端棋盘截图和一张停机截图，均已归档至 [本批次](evidence/a12-cursor-manual-20261111-01/)。两端显示 Cooperative，最终 Alice 在 (2,0)、Bob 在 (0,2)，对方窗口有对应位置；服务端记录 `mode=1 diff=easy players=2`、9×9 / 10 雷、两人注销至 total=0 及 stopped cleanly。静态截图支持有限位置同步，不宣称连续录像或完整规则通过；客户端退出后提示符、具体按键及退出码未单独采集。停机后实际检查 8081 无监听，两源码与已测试快照哈希一致，未重复测试，见 [清理记录](evidence/a12-cursor-manual-20261111-01/cleanup-and-source-check.log)。
- 实现边界及局限：本任务仅涉及上述入口、测试和本地文档，不改 A1.1 成果、依赖、许可证、共享文件或其他项目。本次测试连接已注销至 total=0，无清理超时；现有 Hub.Run 无停止接口，关闭连接并等待注销不能等同所有 goroutine 完整回收，阻塞 Hub 由测试进程退出清理。本任务不顺带改 Hub 生命周期。

## A1.1 交付（指导下完成）

- 问题：双客户端选择 Last Standing，原版仍以单人模式开局；修复前本人观察为 `mode=1 players=2`，预期为模式 2。
- 原因：调用方已创建模式通道并传给 `NewHub`，但构造时未保存到 `Hub.setModeCh`，字段保持 nil。`ForwardSetMode` 的发送无法立即进行，空 `default` 跳过发送，程序继续运行；nil 字段并不表示模式数值为 0，游戏循环的原始默认模式是 1。
- 用户核心改动：在 `internal/server/hub.go` 保存 `setModeCh: setModeCh`；在 `internal/server/hub_test.go` 新增到达检查，以玩家 1、Last Standing 为输入，检查实际收到的 PlayerID 与 Mode。参数改名练习也保留在构造函数中。
- 设计理由：恢复已有参数到字段的交接，沿用现有通道和调用链；不新增抽象、兼容层或消息可靠性策略。
- 分工：核心修复与代表性断言由用户在指导下编写，AI 负责讲解、只读检查、验证与文档。AI 的 gofmt 经对照仅转换两文件的 CRLF/LF，没有改逻辑。不是用户独立设计整个模块。

### 关键调用链

```text
启动：cmd/server/main.go 的 main
      创建 setModeCh → NewHub 保存字段；同一通道交给 snake.Loop

操作：cmd/client/main.go 的 runSnake（模式选择 / 发消息）
      → internal/server/client.go 的 readPump（MsgSetMode）
      → internal/server/hub.go 的 ForwardSetMode（h.setModeCh <- 消息）
      → internal/snake/loop.go 的 runOnce（接收、检查房主并更新 gameMode）
      → 开局日志 mode=2
```

回归测试只检查其中的 Hub 转发：`newTestHub` 准备容量 4 的目标通道，经公开构造器创建 Hub，调用 `ForwardSetMode`，再从同一目标通道接收并检查内容；无需服务器、Hub.Run 或额外 Sleep。

### 已取得的验证

- 同一完整回归测试在原始实现副本中实际失败（无消息，退出 1），在当前修复版通过（退出 0）；服务器包 5 项通过，全库 37 个顶层测试通过、0 跳过，全库构建退出 0。命令、环境与原始输出统一见 [正式回归记录](EVIDENCE.md#a11-正式回归与全库验证2026-11-11联机待验收)及 [日志目录](evidence/a11-regression-20261111-01/)。
- 实际顺序是用户先补字段、后补完整测试；旧实现使用仓库外原始提交副本加相同测试取得失败对照，不宣称按时间先红后绿。首次测试草稿的 vet 格式报告失败另存，不充作模式缺陷的失败证据。
- 本人修复后联机：[服务端截图](evidence/a11-snake-after-fix-20261111-01/server-mode-2-and-clean-shutdown.png)确认 `mode=2 players=2`，两连接注销至 total=0，服务端 stopped cleanly。用户明确确认两端出现棋盘、各按 Q 返回各自提示符；客户端侧属于用户报告，未另采集画面或退出码。
- 停机后实查 8080 无监听，源码与已测试快照一致，未重复测试。联机证据及清理记录见 [EVIDENCE.md](EVIDENCE.md#a11-本人修复后联机2026-11-11工程验证完成)。
- 本次只验证模式转发、选择与双人开局和退出；完整 Last Standing 规则、重连、多轮模式保留、Race、覆盖率、压力、性能及跨机器联机未验证。没有性能提升指标。

## 固定学习结果的实际范围

| 学习项 | 已讲解 / 指导下完成 | 自行修改或判断的实际证据 | 仍需练习的范围 |
| --- | --- | --- | --- |
| 参数、字段、初始化 | 普通整数示例、C# 对照、正常字段交接；两次参数改名练习 | 根据目标自行补齐 `startCh: incomingStartCh`，实际文件及编译核验；能说明后续读取不改名，因为值已保存在字段 | 不据此宣称能独立设计构造函数或维护完整模块 |
| channel 创建、保存、收发 | 正常缓冲收发、无缓冲及 nil 死锁对照，真实字段修复 | 经定位提示补真实字段赋值；根据目标自行创建容量 1 的通道并运行成功 | 继续区分保存通道值、发送一条消息、接收消息；曾把发送 90 误认为初始化通道 |
| select/default | nil/default 与正常缓冲通道演示；补讲发送方向、通信能否立即进行及其与 C# switch 的区别 | 自行修改初始化并预测输出，实测成功；最终对“容量 2、已有 1 条消息”的发送回答“有空位，能发送”，局部边界理由正确 | 最终题为口头预测，未新增运行；此前满通道判断及为何没卡住有混淆，仍需在真实消息流程中巩固 |
| 回归测试 | 输入、动作、到达和内容检查，Errorf / Fatalf；旧失败、新通过的版本对照 | 根据整数示例自行写出两字段判断，实际源码及运行已核验；格式报告方法经提示修正 | 尚未验证从空白独立搭建测试框架、完整解释并发模块或独立调试整个调用链 |

本轮学习收尾已完成，记录为“指导下完成修复、断言及相近小变化”，不升级为全面独立掌握。详细过程、负例和最终口头检查见 [EVIDENCE.md](EVIDENCE.md)。

## 此前基线与本人复验

- 原始基线记录 Go 1.27.3 / Windows amd64 的依赖校验、测试和构建成功；服务端与双客户端进行过有限冒烟。本轮自动验证另列于 A1.1，不混为同批次。
- 原始 Snake Last Standing 缺陷已由 A1.1 修复；历史失败结果仍保留在 BASELINE 与修复前证据中。
- 扫雷 Cooperative / Easy 双端一次翻格同步及退出已在指导下复验；同步、客户端退出与后续服务端停机为分段证据，不合并成一次连续会话。详见 [本人复验](EVIDENCE.md#本人最小复验记录)。
- 编辑器准备已完成：此前 AI 实测 gopls v0.24.0，用户报告 F12 可跳定义。环境排错不计为业务功能贡献。

## 已知风险（A1.1 已处理 R1；A1.2 仅处理 R4 的光标身份归属）

| 编号 | 事实 / 证据 | 影响与状态 |
| --- | --- | --- |
| R1 | 原始字段漏赋值使发送跳过；修复、同测试旧/新对照及本人 `mode=2 players=2` 开局均有证据。 | 本任务修复与有限冒烟完成，完整游戏规则未验证。 |
| R2 | 原始 README / CI 用 Go 1.21，go.mod 要求 1.26.3；README 有 29 tests 徽章、36 tests 示例，原始源码静态计数为 36。 | 配置及文档未修；本轮本机 37 PASS / 0 SKIP 不代表远程 CI 成功。 |
| R3 | 客户端只声明 -name，连接地址固定 localhost:8080/8081；此前已核对 -h。 | README 的 -server 局域网步骤不可用，跨机器未验证。 |
| R4 | WebSocket 允许任意 Origin；原扫雷光标入口采用消息内玩家 ID，读连接限制与 Snake 不同。 | 光标身份来源已由用户修复，同测试失败→通过、全库验证和正常位置同步的有限截图有证据；Origin、读连接限制及其他治理不在本任务范围。 |
| R5 | 存在忽略 WriteMessage 错误、recover/default 丢弃消息及部分测试依赖 Sleep 的路径。 | 本次未改变既有非阻塞策略；通道缓冲满时仍可能跳过转发。可靠性需专项验证。 |
| R6 | 此前 GitHub 直连重置、Go 代理超时，借助本机代理完成下载。 | 网络可变，旧代理配置不代表以后仍可用。 |

## Git 与个人仓库（2026-11-11）

- 用户已授权创建个人副本并使用本人 GitHub 仓库作为远程地址；AI 已通过登录为 LearnedYet 的 GitHub CLI 创建公开 Fork：[LearnedYet/multiplayer-terminal-games](https://github.com/LearnedYet/multiplayer-terminal-games)。GitHub API 确认 parent 为 siyad01/multiplayer-terminal-games。
- 本地仓库原已存在于 `C:\Dev\Go\A\.git`，未重新初始化。`origin` 已指向个人 Fork，原作者地址保留为 `upstream`；本仓库的默认推送远程设为 origin。HEAD 仍为原始提交，尚未提交或上传本地修改。
- 本次只创建上游副本、调整本地 Git 配置并记录事实，不修改 A1.2 业务代码或验收状态，也未进行实际 push。操作、首次参数错误及最终核对见 [Git 配置证据](evidence/git-personal-fork-20261111-01/setup.log)；这项配置不计为个人游戏功能实现。

## 下一步与未授权方向

A1.1、A1.2 均已交付并停止，没有当前阻碍。后续任务须由用户明确开启，再明确范围、核心学习点、作者及验收终点。旧 PORTFOLIO 不再需要同步，现有共同学习规则仍有效。

除已批准的光标身份归属外，其他外部输入治理、连接错误处理、复现性及量化基线、对局生命周期、房间隔离等均未获本任务授权。不自动进入其他功能、提交、推送或部署。
