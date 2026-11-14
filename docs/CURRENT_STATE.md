# 当前状态

## 阶段与边界

- A2 单对局生命周期于 2026-11-11（America/Los_Angeles）开启，2026-11-14 完成并停止。初期代码由用户在 Gemini 指导下编写，2026-11-12 交回本聊天按 TODO 接续；固定名单、连接身份、离开保留成绩、取消、结果快照及进程内保存、递增编号已完成必要工程验证和指导下学习复盘。验收批次全库 43 项测试及构建通过。阶段进度 5e47ebf 及本次接续交付 e3413e3 均已推送个人 origin/main，回执文档另行归档；不自动开启新阶段或提交后续工作。
- A1.2 扫雷光标消息身份归属修复于 2026-11-11 完成。用户指导下保存核心断言及 `c.ID()` 修复，同一测试修复前失败、修复后通过，全库 38 项测试及构建通过，双客户端正常位置同步和停机有有限截图证据；学习复盘已完成。
- A1.1 模式切换修复与回归测试于 2026-11-11 完成交付；用户于 2026-11-02 开启任务，原任务书为 `C:\Dev\Go\docs\A_PHASE1_TASK.md`。其成果及证据保留，不重做验收。
- Phase 0 原始环境与基线于 2026-10-29 记录，本人最小复验于 2026-11-02 收尾；历史结果不重写。
- 原始上游基线：`488dd313bc13de387d12a24f5d05a5c28bed5899`；A2 最终代码交付提交为 `e3413e3b66deaac9113f4ab1e94cd77dcbe59b18`，包含 A1 及此前 A2 进度。2026-11-14 普通推送后已联网 ls-remote 核对个人 GitHub main 同号，见[本次回执](evidence/a2-push-20261114-01/verification.log)；最新回执文档提交以 git log/origin/main 为准，不在文件内自引用自身提交号。原始基线不重写。
- Snake、Minesweeper、WebSocket Hub、二进制协议及终端客户端均为上游已有能力。个人改动、AI 辅助及原始验证分别记录于 [EVIDENCE.md](EVIDENCE.md)。
- A1.1/A1.2 和 A2 本轮均未改依赖、许可证或共享文件；此次已授权的提交和个人仓库推送均完成，未部署。构建产物与临时副本仍放仓库外。
- 用户说明 PORTFOLIO.md 为旧 Gemini 同步材料，已不再作为当前进度依据或同步待办。当前进度以本项目状态和证据为准，共同学习规则继续遵守。

## A2 阶段收尾（2026-11-14 完成停止）

- 最新结论（2026-11-14）：A2 必要工程验证及指导下学习复盘已完成，阶段完成停止，不新增练习或自动开启后续阶段。验收批次全库 `go test -mod=readonly ./... -count=1 -v -timeout=60s` 为 43 PASS、0 SKIP、0 FAIL，`go build -mod=readonly ./...` 退出 0；收尾时只读核对 cmd/internal 的 21 个 Go 文件与验收哈希差异为 0，没有重复运行测试，见 [验收批次](evidence/a2-acceptance-20261114-01/)。以下逐步记录保留历史来源，最新状态以本条为准。
- 连接与生命周期有限诊断：真实 Hub 编号 0 被复用但 session 从 1 递增到 2，旧/新 rejoin 连续请求只转发当前连接；慢连接队列饱和后收到其 ID/SessionID 的 leave，随后仍能接纳新连接。本机 httptest + 真实 WebSocket + mineWSHandler/Hub/Loop 完成 Alice/Bob 开局、Alice 断开、Carol 复用 0 被拒绝且插旗未改棋盘、Bob 有效插旗后断开取消、Carol 原连接再加入后胜局及败局、最后 Dave 进入新大厅；Loop 日志显示 mine - 1/2/3 及保存数量 1/2/3。
- 诊断边界：上述为 AI 输出观察值的临时 overlay，未代写核心业务断言；胜负路径通过临时工厂 hook 取得 Game 指针选择安全格/雷格，未更改 RNG 或规则。结束标记/成绩来自 Game.SnapshotResult 的观察，Loop 保存次数来自原生产日志，没有读取未暴露的历史列表或冒充查询接口。正式编号复用回归及快照独立性由用户断言另行支持。未执行 Race、性能或手动终端画面验收。结果仅进程内保存，重启会丢失，编号只保证本次 Loop 内唯一。
- 学习复盘完成（2026-11-14）：用户能说明连接身份保护原玩家记录、未持久化导致重启丢失。初次对拒绝 join 后仍检查事件表示不明白，经实际提前取消回归场景解释，用户回答“拒绝加入复用了编号，也被监测，会产生消息，所以要加身份识别”，抓住拒绝加入后连接仍产生日志/事件、需核对身份的核心。AI 补正先分配编号后拒绝本局接纳的顺序；记录为在提示和解释下完成关键概念说明，不升级为独立设计或熟练掌握整个模块。无新增代码练习或测试。

- 接续调查（2026-11-13）：本轮已读取项目/共享规则、状态及证据；HEAD 与本地 origin/main 同为 5e47ebf，未联网重查远程。调查前未提交内容仅为三份交接文档及 a2-handoff 回执；`git diff --exit-code HEAD -- cmd internal go.mod go.sum LICENSE` 返回 0。本轮未重跑 Go 测试或构建。
- 固定名单链路实查：mineWSHandler 调用 Hub.Register 后即启动读消息；Hub.assignID 从在线 clients 中选择最小空闲编号，unregister 删除后可立即复用。调查时 GAME join 仍会 AddPlayer（下条已修正）；即使回复 false，Hub 也只发错误消息，读循环仍会把该连接的 c.ID() 操作送入 mineCh。leaveCh、joinCh、mineCh 分属不同通道，不能假定循环先处理旧玩家 leave 再处理复用编号的新连接输入；Left 检查不足以证明该窗口安全。以上为静态路径分析，尚无本轮运行复现。
- GAME join 小步已保存并审查（2026-11-13）：用户移除中途 AddPlayer 和棋盘广播，只通过 ResultChan 回复 false，不退出本局；大厅 join 未改。扫雷包现有 13 项测试通过、退出 0，源码前后哈希一致，见 [本批次](evidence/a2-fixed-roster-join-20261113-01/package-tests.log)。现有测试未覆盖 runOnce 的 join，仅支持既有行为回归，不据此认定连接身份隔离完成。
- 连接身份方案已批准（2026-11-13）：用户同意保留现有连接等待/再加入行为，在服务端为每个连接增加独立身份，并随接纳、输入、离开等相关事件传递，让本局核对开局时的连接。必要的内部消息调整属于本方案，不改变客户端线协议、不引入账号或跨重启身份；核心代码及断言仍由用户按 TODO 填写，同一方案内不逐步重复请示。
- 连接身份分配已保存并审查（2026-11-13）：用户在 Client 保存 sessionID（uint64），Hub 保存 sessionCounter（sync/atomic.Uint64），NewClient 用 Add(1) 分配身份，SessionID 方法只返回保存值。身份在注册及读消息启动前固定；HTTP 可并发创建连接，计数器原子递增。AI 未改源码，server 包既有 5 项测试通过、退出 0，源码前后哈希一致，见 [分配验证批次](evidence/a2-session-allocation-20261113-01/package-tests.log)。现有测试未断言新身份唯一性或同连接稳定性，不能据此宣称身份隔离通过。
- 消息传递已保存并审查（2026-11-13）：用户补齐 JoinMsg、LeaveMsg 的 SessionID uint64，连同已保存的 MineInputMsg、扫雷 SetModeMsg，相关发送处均从当前连接取身份。上一草稿三处 unknown field 的编译失败保留于 [失败批次](evidence/a2-session-messages-draft-20261113-01/package-tests.log)，本轮已消除。Snake loop 的其余变化经移除两个新增字段后与格式化 HEAD 比较相同，属格式整理，无额外逻辑变化。
- 本轮验证：`go test -mod=readonly ./internal/server ./internal/snake ./internal/minesweeper ./cmd/mineserver -count=1 -v -timeout=30s` 为 30 PASS、0 SKIP、0 FAIL、退出 0，源码前后哈希一致，见 [消息传递批次](evidence/a2-session-messages-20261113-01/package-tests.log)。既有真实 WebSocket 测试输出收到 PlayerID=0、SessionID=1 的光标消息，但未新增 session 核心断言，不认定固定名单/复用身份隔离通过；未执行全库或单独构建。
- 开局身份传递已保存并审查（2026-11-13）：用户为 StartMsg 增加 SessionID、ForwardStart 增加 sessionID uint64 参数并填写消息，扫雷 handler 和共享 client.readPump 三处调用均传入发送连接身份；Snake 判定未改，client.go 其余差异为格式整理。相关四包相同命令再次 30 PASS、0 SKIP、0 FAIL、退出 0，源码前后哈希一致，见 [开局传递批次](evidence/a2-session-start-20261113-01/package-tests.log)。现有测试未断言 StartMsg 的 SessionID，不能据此认定开局身份隔离已验证。
- 大厅身份匹配已保存并静态审查（2026-11-13）：用户创建 playerSessions，成功 join 时保存身份及名字，leave/setMode/start 均用 map 双返回值判断存在且身份相同；有效 leave 删除两张表，原房主判断保留。用户已修正为先更新身份与名字再 broadcastLobby，AI 未代改。扫雷包既有 13 项测试通过、退出 0，源码前后哈希一致，见 [大厅批次](evidence/a2-session-lobby-20261113-01/package-tests.log)；既有用例未调用 runOnce，不将包通过当作大厅身份场景已运行验证。
- GAME 身份匹配已保存并审查（2026-11-13）：用户在 mineCh 和 leaveCh 分支先检查编号存在且 SessionID 匹配，不符合就 continue；有效 leave 保留 RemovePlayer、取消和快照逻辑，playerSessions 在 GAME 中未修改。AI 未代写业务，新增测试外壳前既有扫雷 13 项测试通过、退出 0，源码前后哈希一致，见 [GAME 批次](evidence/a2-session-game-20261113-01/package-tests.log)。
- 编号复用核心断言已保存并通过（2026-11-14）：用户填写拒绝中途加入、不得提前结束、结果非空/取消/两人名单，并按名字查找 Alice/Bob、确认 Alice.Flags=1，删除占位失败；AI 读取实际测试，未代写断言。相同测试对照 a2-session-lobby 的旧 loop 快照：endedEarly=true，在用户“对局在离开前结束了”断言失败、退出 1；当前实现 endedEarly=false、返回两人取消结果、Alice.Flags=1，1 PASS、退出 0。测试字节及相关源码前后哈希一致，见 [本批次](evidence/a2-session-regression-20261114-01/)。这是内部通道场景，未包含真实 Hub 编号分配/网络断开。
- Hub 收尾已保存并验证（2026-11-14）：用户移除重复存在性检查，保留编号与 Client 指针检查，慢连接移除后发送带正确身份的 LeaveMsg，通知移至关闭函数调用后。AI 已说明通知位置属于代码组织取舍，正常路径放在匿名函数内并非功能错误，不应将位置整理强制当作验证前置条件；用户后来已自行整理完成。实际诊断与全库结果见本节最新结论。

- 阶段目标：沿用现有单房间 Cooperative / Easy 玩法，为每局保存稳定编号及本局参与者记录，明确开始、胜负、取消和再次开局，并形成结束后不再随游戏状态变化的结果。本阶段不接入 MySQL/Redis、多房间或账号系统。
- 核心学习点：在线连接名单与本局参与记录的不同用途；状态转换、数据保存期限、结果快照及回归断言。既有核心代码和断言由用户在 Gemini 指导下编写，本聊天现已接回带学、审查及验证；后续仍由用户编辑核心内容，不将其他指导工具的工作记为本聊天实现。
- 阶段开启时的基准现状：`Loop` 反复调用 `runOnce`；大厅接到房主 start 后创建 Game 并加入玩家；游戏中允许 Cooperative 中途加入；原版 leave 调用 `Game.RemovePlayer` 删除含 Reveals/Flags 的玩家对象及光标。胜负时广播后返回下一轮大厅；全员离开分支未见取消结果。上述为源码调查，未执行新的场景测试。
- 第一个问题示例：Alice 本局已有翻格贡献，离开后仍应能在本局结果里找到；原版 RemovePlayer 会删除该记录。当前改动已使用 Left 标记离开并保留玩家对象，用户已按指导补强“离开前后成绩完全相同”的断言，定向测试通过；后续接续操作限制、循环与结束结果。
- 已确认首版业务规则（2026-11-12）：开局固定本局名单，新连接不加入当前局；离开者停止参与操作但保留本局成绩；胜负沿用当前 Cooperative 规则；全员离开把本局记为取消；结束结果固定，再次开局产生新编号。此前保留中途加入的备选本次不采用。
- 实现前注意：Hub 会复用短玩家编号，不能把保留下来的本局参与记录简单等同新的在线连接；拒绝中途加入也需在输入接纳与连接管理链路验证，不能只改 join 回复后认定安全。编号形式、具体结果结构及规则实现待逐步说明，不预先批准新依赖或通用抽象。
- 验收方向：结束后不再改分，离开者本局成绩保留，全员离开能取消，连续两局编号和结果不混用；加入政策按确认规则验证，保留相关测试/构建及有限联机证据，并按实际表现记录学习范围。当前通过的是早期模型用例及扫雷包测试，不代表 A2 验收完成。
- 当前文件实查（2026-11-13）：用户已在 Loop 维护递增 matchNumber，生成 matchID 后沿 runOnce → NewGame 保存至 Game.MatchId，移除 Game 中的时间和随机编号表达式，测试辅助构造器提供固定编号。与修复前相同的两局诊断现显示 mine - 1、mine - 2，保存数量 1、2；取消快照正确、随后大厅可接受新玩家。Cooperative 中途加入仍未处理，A2 未完成。
- 接回验证：`go test -mod=readonly ./internal/minesweeper -run '^TestPlayerLeaveRetainsScore$' -count=1 -v -timeout=30s` 通过；`go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s` 为 10 PASS、0 SKIP、0 FAIL，均退出 0。Go 1.27.3 / Windows amd64，进程 GOWORK=off；源码前后哈希一致。日志及原始快照见 [接回检查批次](evidence/a2-handoff-check-20261112-01/)。本轮未执行全库测试、构建、Race 或联机，不沿用 A1 的 38 PASS 证明当前 A2。
- 最新小步（2026-11-12）：用户已保存 `g.Players[0].Reveals != scoreBeforeLeave` 判断，Fatalf 参数按“离开后、离开前”与提示文本一致。AI 读取实际文件后执行同一定向命令，1 PASS、退出 0，源码前后哈希一致；本步只改断言，未重复包测试或扩大验证。见[成绩断言验证](evidence/a2-score-assertion-20261112-01/targeted-test.log)，不记为独立设计测试或业务修复。
- 插旗小步已保存（2026-11-12）：用户已将原插旗检查移至独立 TestLeftPlayerCannotFlag，使用未翻开的棋盘，Alice 离开后对 (0,0) 插旗，检查返回 changes 为空且目标仍隐藏。AI 读取实际文件后运行 `go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s`：11 PASS、0 SKIP、0 FAIL，退出 0；源码前后哈希一致，见[本批次日志](evidence/a2-left-flag-20261112-01/package-tests.log)。代码来自此前已提供的完整示例，仍记为指导下完成；当前不再使用完整答案带学。
- 光标练习完成（2026-11-12）：用户按思路及 TODO 骨架自行补齐 TestLeftPlayerCannotMoveCursor 的动作与三次数量断言，分别检查加入后一个光标、离开后无光标、离开后尝试移动仍无光标；每次动作后重新查询 AllCursors，删除了占位失败。AI 读取实际文件并运行 `go test -mod=readonly ./internal/minesweeper -run '^TestLeftPlayerCannotMoveCursor$' -count=1 -v -timeout=30s`：1 PASS、退出 0，源码前后哈希一致，见[本批次](evidence/a2-left-cursor-20261112-01/targeted-test.log)。仅认定完成给定骨架的相似小变化，不升级为独立设计整个模块。
- 取消路径核验（2026-11-13）：用户第三次保存加入 Over 提前返回，并将 GAME leave 改为 return false。AI 临时 overlay 诊断（不改源码/核心断言）观察到进行中的 Game 取消后 Over=true、Canceled=true；胜局/败局调用 Cancel 后保留原结果，重复调用不变；Cooperative / Easy 两人局中 Alice 离开后 Bob 光标广播仍可收到，Bob 离开后出现 restarting / lobby open 并接受 Carol。诊断执行成功退出 0，业务/测试源码前后哈希一致，见[诊断原始日志](evidence/a2-cancel-observe-20261113-01/diagnostic.log)。这是实际调用的有限诊断，不是新增业务断言回归测试或双客户端联机，固定结果快照尚未实现。
- 结果类型已保存（2026-11-13）：用户按 TODO 定义 MatchResult，实际字段为 GameID string、Won bool、Cancel bool、PlayerResult []MinePlayer，字段类型符合当前需要，保留用户命名。AI 只读核对，未修改源码；本步仅声明数据类型，未新增程序测试，不据此宣称结果已复制、冻结或保存。
- 快照方法已保存（2026-11-13）：用户按 TODO 补齐 SnapshotResult，res 复制 GameID/Won/Cancel，新建长度 0、容量 len(g.Players) 的 []MinePlayer，遍历 append 玩家值 *p，再赋给 res.PlayerResult 并返回。AI 静态核对正确，包测试 12 PASS、0 SKIP、0 FAIL、退出 0，源码前后哈希一致，见[本批次日志](evidence/a2-snapshot-method-20261113-01/package-tests.log)。既有 12 项未调用 SnapshotResult，快照独立性尚待业务用例验证；结束流程尚未调用和保存结果。
- 快照独立性已验证（2026-11-13）：用户按 TODO 补齐 TestSnapshotResultCopiesPlayers：Alice 7 分取消后取快照，将原玩家改为 99 分并 RemovePlayer，断言快照分数仍为 7、Left 仍为 false。定向运行 1 PASS、退出 0，源码前后哈希一致，见[本批次日志](evidence/a2-snapshot-isolation-20261113-01/targeted-test.log)。分数失败提示目前误读 g.Players 的数值而非快照值，提示用户顺手修正；断言条件本身正确。
- 保存路径核验（2026-11-13）：用户已修正取消快照顺序和失败提示。包测试 13 PASS、0 SKIP、0 FAIL、退出 0；AI 仓库外 overlay 诊断观察到 runOnce 返回 stopped=false、Cancel=true 且两位玩家 Left=true；真实 Loop 两轮取消后日志数量 1、2，第三个大厅接受新玩家。但两轮编号均为 mine-1794578159024-4825，编号唯一性实际失败。见[本批次诊断](evidence/a2-result-save-20261113-01/diagnostic.log)及[包测试](evidence/a2-result-save-20261113-01/package-tests.log)。诊断仅输出观察值，没有用户核心断言，不能把诊断 PASS 当作编号验收通过；胜负返回结果仍仅静态核对。结果仅进程内保留，尚无外部历史查询接口。
- 本次收尾：编号碰撞已在单 Loop 进程内修复，用户按 TODO 保存三个文件，AI 未代写业务。重复诊断正常退出 0；`go test -mod=readonly ./... -count=1 -v -timeout=60s` 为 42 PASS、0 SKIP、0 FAIL、退出 0，`go build -mod=readonly ./...` 退出 0，源码前后哈希一致，见[编号修复批次](evidence/a2-id-fix-20261113-01/)。提交推送已完成；交接时再次只读核对源码与该批次哈希相同，没有重跑测试。
- 新聊天第一步：读取项目及共享规则、本状态和 EVIDENCE，核对当前差异，然后追踪 cmd/mineserver 的连接注册/读消息路径 → Hub 的加入接纳、ID 分配/复用 → runOnce 的 GAME join 分支。目标是落实已确认的固定开局名单，阻止新连接加入或借复用编号影响本局；先读清实际链路再给最小 TODO，不只改 join 回复便宣称完成。必要相邻变更须说明理由，不扩展其他网络治理。
- 交接分工：核心业务及核心断言仅给思路、准确位置和 TODO 骨架，用户保存后再审查；本轮用户已按骨架补齐光标、快照隔离、取消与编号相关小变化，不认定完全独立设计模块。不要恢复此前暂停的相似翻格练习；只做与剩余业务相关的必要验证。
- 工作区说明：推送完成时工作区干净；随后本次换聊天仅补充 AGENTS、CURRENT_STATE、EVIDENCE 和 a2-handoff-20261113-01 回执批次，尚未提交。这些是本聊天交接文档，无待审业务改动；此前提交/推送授权不自动延续到新任务。
- 带学方式更新（2026-11-12）：用户要求后续业务代码及核心断言只给思路和 TODO 骨架，不贴完整答案，由其自行填写。此前报告“测试跑通了”时新插旗用例尚未保存，未据此认定通过；现已保存并由本聊天实测通过，按上条更新。AI 仅提供提示、审查和验证，不代填下一小步核心断言。
- 编号边界：现有毫秒时间加游戏 rng 的四位随机数已经在相同批次复现连续两局重号；本步递增编号仅保证同一 Loop 运行期间的唯一性，进程重启后可重复，与本阶段进程内结果列表的寿命一致。跨进程、跨重启标识及磁盘持久性不在当前实现中。

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
| R7 | A2 读码时看到 Game.Reveal 直接以坐标计算 Cells 索引，当前扫雷入口仅检查消息长度，未见棋盘坐标范围检查。 | 仅静态记录，未运行异常坐标用例；独立输入治理问题，不在本次生命周期开阶段时顺带修复。 |

## Git 与个人仓库（2026-11-11）

- 用户已授权创建个人副本并使用本人 GitHub 仓库作为远程地址；AI 已通过登录为 LearnedYet 的 GitHub CLI 创建公开 Fork：[LearnedYet/multiplayer-terminal-games](https://github.com/LearnedYet/multiplayer-terminal-games)。GitHub API 确认 parent 为 siyad01/multiplayer-terminal-games。
- 本地仓库原已存在于 `C:\Dev\Go\A\.git`，未重新初始化。`origin` 指向个人 Fork，原作者地址保留为 `upstream`；默认推送远程为 origin。用户本轮授权本地提交，HEAD 已为 `d02f621`，未推送；详见 [提交与启动记录](evidence/a2-start-20261111-01/local-commit-and-phase-start.log)。
- 本次只创建上游副本、调整本地 Git 配置并记录事实，不修改 A1.2 业务代码或验收状态，也未进行实际 push。操作、首次参数错误及最终核对见 [Git 配置证据](evidence/git-personal-fork-20261111-01/setup.log)；这项配置不计为个人游戏功能实现。

## 下一步与未授权方向

A1.1/A1.2 历史成果保留；A2 已接通取消、快照复制、循环内留存和进程内递增编号，42 项全库测试与构建通过，进度已提交推送为 5e47ebf。A2 尚未完整交付；新聊天接续固定开局名单及必要综合验收，不另开阶段。保留 Gemini 与本聊天的指导归属，继续只给用户思路和 TODO 骨架。旧 PORTFOLIO 不再需要同步。

其他输入治理、量化基线、房间隔离及数据库均不在当前 A2 范围，不自动推送、部署或安装组件。
