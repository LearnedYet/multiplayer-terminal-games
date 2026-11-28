# 当前状态

## 阶段与边界

- A3 MySQL 结算与查询于 2026-11-15 获用户明确同意开启，2026-11-23 完成本机工程验收。范围为单房间 Cooperative / Easy 的事务结算、防重复提交、按对局查询、有限分页和浏览器战绩入口。用户在指导下完成 SQL、结算及查询核心代码，AI 负责审查、验证、浏览器试玩实现和文档。2026-11-24 用户确认已推送 `36e97a5` 至个人 `origin/main`，本轮核对本地 `main` 与 `origin/main` 跟踪引用均指向该提交。[A3 验收证据](evidence/a3-acceptance-20261123-01/verification.log)
- A2 单对局生命周期于 2026-11-11（America/Los_Angeles）开启，2026-11-14 完成并停止。初期代码由用户在 Gemini 指导下编写，2026-11-12 交回本聊天按 TODO 接续；固定名单、连接身份、离开保留成绩、取消、结果快照及进程内保存、递增编号已完成必要工程验证和指导下学习复盘。验收批次全库 43 项测试及构建通过。阶段进度 5e47ebf 及本次接续交付 e3413e3 均已推送个人 origin/main，回执文档另行归档；不自动开启新阶段或提交后续工作。
- A1.2 扫雷光标消息身份归属修复于 2026-11-11 完成。用户指导下保存核心断言及 `c.ID()` 修复，同一测试修复前失败、修复后通过，全库 38 项测试及构建通过，双客户端正常位置同步和停机有有限截图证据；学习复盘已完成。
- A1.1 模式切换修复与回归测试于 2026-11-11 完成交付；用户于 2026-11-02 开启任务，原任务书为 `C:\Dev\Go\docs\A_PHASE1_TASK.md`。其成果及证据保留，不重做验收。
- Phase 0 原始环境与基线于 2026-10-29 记录，本人最小复验于 2026-11-02 收尾；历史结果不重写。
- 原始上游基线：`488dd313bc13de387d12a24f5d05a5c28bed5899`；A2 最终代码交付提交为 `e3413e3b66deaac9113f4ab1e94cd77dcbe59b18`，包含 A1 及此前 A2 进度。2026-11-14 普通推送后已联网 ls-remote 核对个人 GitHub main 同号，见[本次回执](evidence/a2-push-20261114-01/verification.log)；最新回执文档提交以 git log/origin/main 为准，不在文件内自引用自身提交号。原始基线不重写。
- Snake、Minesweeper、WebSocket Hub、二进制协议及终端客户端均为上游已有能力。个人改动、AI 辅助及原始验证分别记录于 [EVIDENCE.md](EVIDENCE.md)。
- A1.1/A1.2 和 A2 本轮均未改依赖、许可证或共享文件；此次已授权的提交和个人仓库推送均完成，未部署。构建产物与临时副本仍放仓库外。
- 用户说明 PORTFOLIO.md 为旧 Gemini 同步材料，已不再作为当前进度依据或同步待办。当前进度以本项目状态和证据为准，共同学习规则继续遵守。

## A5 双房间隔离（2026-11-28 本机验收完成）

- 已按 [总路线图](../../docs/ROADMAP.md#L45) 于 2026-11-28 开启并完成 A5 本机验收，项目 A 首版（单玩法生命周期、MySQL 事务结算与查询、固定双房间隔离）收尾：
  - **穿帮复现与双房间核心代码**：AI 在试玩前端（`tools/playtest/public/`）加入 `1 号房 / 2 号房` 选择器并在未改后端的 `8093` 复现两房挤入同一大厅。随后 AI 将 `cmd/mineserver/main.go` 单房通道样板提取为 `mineRoom` / `startMineRoom`，用户在指导下完成 `rooms` 固定双房映射、`mineRoomWSHandler` 读取 `?room=`（空默认 `"1"`、未知返回 `400`）、`/ws` 路由接线以及停机遍历 `<-room.done`。
  - **4 客户端实机联调与入库**：在 `8093` 与本机 MySQL 8.4.13 上由用户亲手操作 4 个浏览器客户端分 1/2 号房（各 2 人）同时开局；房间 2 踩雷结算入库（`5VG2HH6B4F35FCRM2TVYPNZ2Q2`，`搭档2` 翻 46 格）并返回大厅期间，房间 1 保持独立运行不受影响。
  - **自动化测试与构建**：用户在 `cmd/mineserver/main_test.go` 的 `TestMineRoomRoutingIsolatesRooms` 编写未知房间 99 返回 400 的状态码断言，AI 辅助完成 4 客户端 WebSocket 外壳、跨房隔离与停机退出集成断言（用户明确后续此类网络测试机器人外壳由 AI 负责，用户负责核心规则小测试、预期边界与前端手动验收）。全库 `45` 项顶层测试全部 `PASS`（0 SKIP、0 FAIL），`go build -mod=readonly ./...` 退出 `0`，见 [A5 验收证据](evidence/a5-acceptance-20261128-01/verification.log)。本次改动尚未提交或推送。

## 原路线 A4：单房间交付点（2026-11-27 本机验收完成）

- [历史路线](../../docs/ROADMAP.md#L44)中的 A4（单房间中间交付点）已于 2026-11-27 完成本机验收：
  - **完整路径与重启恢复**：用户与 AI 浏览器搭档在 `8093` 试玩联机（用户房主 + AI 搭档），AI 搭档插旗 1、翻格 1 后中途离开，用户继续通关（`YGXQJO6XJI7VYGQQCCO4N5VS2E`，70 + 1 = 71 安全格）。MySQL 实查两表一致；随后停掉旧 `8093` 进程并启动新进程，战绩列表及单局详情刷新仍可查。
  - **本人在指导下修复小规则（R8）**：先在 `8093` 复现“玩家B 插旗、玩家A 拔掉后，战绩详情里玩家B 仍为 1”；用户在指导下完成 `internal/minesweeper/game.go`（`Flag` 的 `case StateFlagged:` 先保存 `flaggedBy := cell.FlaggedBy` 再清零并按 `flaggedBy` 扣减，首次误写为 `playerID` 后自行修正）及 `internal/minesweeper/game_test.go` 的 `TestCrossPlayerUnflagDecrementsOriginalOwner`。重新构建 `8093` 实测双方插旗数均归 `0`，且用户能自己推演交叉插拔旗与 `uint16` 防下溢边界。
  - **断库故障与恢复观察**：运行 `8093` 期间正常停止本机 MySQL，`JTHBPLKWQCAXTGRUBEQSRYOQFG` 结束时记录写库失败（`connectex` 连接被拒）后继续返回大厅开放下一局，服务未崩溃；重新拉起 MySQL（13306，库内 8 局）后无需重启 `8093` 即恢复战绩查询。
  - **全库验证**：`go test -mod=readonly ./... -count=1 -v -timeout=60s` 44 项顶层测试全部通过（0 SKIP、0 FAIL），`go build -mod=readonly ./...` 退出 0，见 [A4 验收证据](evidence/a4-acceptance-20261127-01/verification.log)。本次改动及证据随后提交并推送至个人 `origin/main`（具体提交号以 `git log` 为准）。

- **2026-11-28 交付后复核**：规则修复为 `fc8c71c`，README 为 `94951ca`。只读查询确认个人 GitHub `main` 为 `94951ca`。用户的同一回归测试在 `36e97a5` 的旧 `game.go` 上因双方计数 `0,1` 失败，在当前提交通过。全库 44 项通过，构建通过。该测试单独运行时，插旗完全失效也会通过，可补上撤旗前确认原玩家已有一面旗的断言。README 待核正两处表述：写库仍同步等结果或 3 秒超时才重新开大厅，规则修复按用户在 Gemini 指导下完成记录。[复核日志](evidence/a4-review-20261128-01/verification.log)

## 账号与「我的战绩」候选讨论（2026-11-24，待确认）

- A3 收尾聊天把注册、登录、固定账号 ID 和「我的战绩」作为可能的后续方向，沿用现有「雷区搭档」页面做可点验收。2026-11-27 的 A4 交付按上节的原单房间范围完成，账号方向留待后续确认。建议限定在本机浏览器试玩的单房间 Cooperative / Easy，保留游客试玩和全部战绩入口。账号仅从登录后新加入的对局开始累计战绩，A3 历史行和游客行的账号关联留空。
- 候选数据方案为 MySQL 分配账号主键，单独保存账号资料，向 `match_players` 增加可空的 `account_id`。服务端从登录会话取得账号 ID，在玩家加入时绑定，在开局冻结名单时保存到本局参与记录，经结算快照写入原事务。「我的战绩」查询按当前登录账号过滤，沿用有限分页和单局详情。现有 `PlayerID` 用于游戏内识别，`SessionID` 用于区分连接实例，账号 ID 跨重连、重新登录和服务重启保持一致。
- 候选登录方案为本机同源网页、带随机盐的密码派生结果、浏览器 Cookie 与服务端内存登录会话。刷新可保持登录，退出或重启服务后需重新登录。首版范围包括注册校验、登录、退出和个人战绩，公开部署、邮箱验证、密码找回、第三方登录、排行榜及历史归属补录留待另议。
- 玩家可见验收分四步。先在注册页创建账号并看到重名提示，再用登录页查看当前账号 ID，刷新、退出和重登后确认身份。接着从大厅以账号身份和游客各打一局，从已有详情查看成绩，辅助核对库里的账号关联。最后从大厅点进「我的战绩」，用两套独立浏览器会话验证两人都能看到共同参与的对局，只有参与者能在自己的列表看到另一局，重启后重新登录仍可查询。每步保留可重复打开的本机入口，使用保存后的实现及真实数据。
- 核心学习点是登录身份如何传到扫雷连接、进入固定名单并关联结算，以及用当前账号范围做查询。若开启账号方向，用户按中文 TODO 填写业务代码、SQL 和代表性断言并保存，AI 负责定位、提示、骨架、审查、试玩界面辅助、验证与文档。

## A3 实施与验收（2026-11-23）

- 已读当前规则、状态和相关真实源码。现有链路为 cmd/mineserver/main.go → minesweeper.Loop → runOnce → Game.SnapshotResult → Loop 内 results 切片；结果仍在进程内，编号现已改用标准库随机文本。当前未接入数据库或增加核心断言。
- 本轮环境：Go 1.27.3 / Windows amd64；PATH 未找到 docker/mysql/mysqld，相关 Windows 服务计数为 0，3306/13306 无监听。wsl --list --verbose 退出 1，输出未安装 WSL。用户确认没有现成 MySQL，要求先提出本机方案。System32 的 vcruntime140.dll/msvcp140.dll 存在，但未启动 MySQL 验证运行条件。见[环境原始记录](evidence/a3-preparation-20261115-01/environment.log)。这些检查不等同于穷举所有安装位置。
- 环境方案已批准并执行：官方 MySQL 8.4.13 Windows ZIP，在 C:\Users\c8618\AppData\Local\GoA\mysql 下保存程序、数据、本地配置及启停/连接辅助脚本，仅绑定 127.0.0.1:13306。未注册服务、修改永久 PATH 或防火墙。独立学习数据库 go_a_learning、账号 a3_dev@127.0.0.1 已创建，业务表数量为 0，尚未接入 Go。root 和学习账号密码仅保存到仓库外 private 下的 client 配置，不复制到证据。包校验、连接、正常停止和重启核验通过，当前实例已停止，见[本批次](evidence/a3-mysql-setup-20261115-01/)。官方依据见 [Windows 安装说明](https://dev.mysql.com/doc/refman/8.4/en/windows-installation.html)和[手动配置步骤](https://dev.mysql.com/doc/refman/8.4/en/windows-install-archive.html)。
- 启停辅助脚本首次重启失败：只等端口消失时旧进程仍释放数据文件，新进程报告 ibdata1 must be writable；修正 Stop 脚本为正常 shutdown 后等待实际 mysqld 进程退出，随后重启通过。Start-Process 返回的启动器 PID 不能替代实际监听进程，已改用 mysqld.pid 并在核验中检查监听归属。AI 负责此环境修正，没有给用户追加测试练习。失败日志和修正后结果均保留。
- 编号改动已由用户按 TODO 保存并审查：Loop 移除 matchNumber 及 fmt，每轮调用 cryptorand.Text 生成一次 matchID，继续传入 runOnce → NewGame → Game.MatchId → SnapshotResult.GameID；结果保存位置没有重新生成编号。既有 math/rand 仍用于棋盘。使用标准库、不新增依赖；随机碰撞概率极低，不宣称绝对唯一。同一局以后再次保存应沿用原编号，未来数据库唯一约束仍待实现；备选数据库分配编号未采用。
- 本小步验证：Go 1.27.3 / Windows amd64，GOWORK=off；go test -mod=readonly ./internal/minesweeper -count=1 -v -timeout=30s 为 14 PASS、0 FAIL、0 SKIP，mineserver 定向构建退出 0，产物在仓库外。21 个 Go 源文件前后哈希差异为 0，AI 未改核心内容。见[本批次](evidence/a3-random-id-20261115-01/)。既有测试未调用生产 Loop，支持相关规则回归，不声称通过了跨进程唯一性或数据库持久化测试；本轮未执行全库、Race、数据库业务或网络验收，未新增随机数重复抽样测试。
- 字段归属讨论已完成：用户列出整局比赛 ID/是否胜利/是否取消，参与者名字/座位号/得分/旗子/是否离开，并指出用 game id 关联本局。与现有结构一致；AI 明确得分对应 Reveals（翻开格数），旗子保存现有 Flags 字段结算值，不改统计规则。该回答记录为根据已有结构完成归类，不升级为独立设计数据库。
- 两表方案已获用户明确同意：match_results 每行一局，game_id 为主键，保存 won/canceled；match_players 每行某局的一名参与者，保存 game_id/player_id/name/reveals/flags/left_early，以 (game_id, player_id) 为联合主键，game_id 外键关联 match_results。两表用 InnoDB，所有结果字段为必填；编号为字符串、玩家编号/成绩计数按现有非负整数范围选型、状态用布尔字段。备选一表 JSON 未采用。查询排序/分页、事务与重复提交内容冲突处理仍待后续设计，驱动依赖变更需明确后再操作。
- 对局表小步已完成：用户在 db/schema.sql 填好 game_id、won、canceled 的类型、非空与主键并保存。AI 审查后把当时完整文件在本机 MySQL 8.4.13 的 go_a_learning 执行，退出 0；SHOW CREATE TABLE 确认三列、game_id 主键及 InnoDB，BOOLEAN 显示为 tinyint(1)。业务表仅 match_results，随后正常停机，13306 无监听。见[原始建表记录](evidence/a3-match-results-20261118-01/verification.log)。
- 参与者表小步已完成：用户按 TODO 填好六列、联合主键和外键。AI 按用户授权移除外键列清单里的多余逗号，整理缩进和文件头说明。仅执行 match_players 的建表语句，客户端退出 0；SHOW CREATE TABLE 确认 (game_id, player_id) 联合主键和 game_id → match_results.game_id 外键。两表均存在，实例正常停止，13306 无监听。见[原始建表记录](evidence/a3-match-players-20261118-01/verification.log)。
- Go 接入方案已由用户确认：标准库 database/sql 加 github.com/go-sql-driver/mysql v1.10.1，不加 ORM；连接参数通过进程环境变量 MINESWEEPER_MYSQL_DSN 传入，不入仓库或日志。AI 已添加驱动及其间接依赖 filippo.io/edwards25519 v1.2.0；Go 版本声明和既有依赖版本无变化。官方校验站直连超时，进程内使用 goproxy.cn 重试并保留签名校验，驱动和 mineserver 定向构建均退出 0。见[依赖准备记录](evidence/a3-mysql-driver-20261118-01/verification.log)。
- 连接探测已完成：用户在 cmd/mineserver/main.go 填写 probeMySQL 和启动分支。配置 DSN 时使用 sql.Open 和带 3 秒超时的 PingContext，结束关闭探测句柄；AI 补上漏掉的驱动注册导入、调整 go.mod 的直接依赖标记并清理旧 TODO。定向构建通过；已配置且 MySQL 运行时启动监听 8081，MySQL 停止时报告连接失败、退出 1 且不监听 8081，未配置且 MySQL 停止时仍正常启动。核验后 MySQL 已正常停止，两个端口均无监听。[原始验证记录](evidence/a3-mysql-probe-20261118-01/verification.log)。
- 已追踪结算链路：胜、负、全员离开取消三个分支均由 runOnce 取 SnapshotResult 并返回非空结果；服务停止分支返回 nil。外层 Loop 的 `if result != nil` 是统一保存入口，目前只 append 到进程内切片。启动探测结束会关闭句柄，不能直接拿来写入。
- 用户已同意结算方案：配置了 MySQL 时由 mineserver 启动时创建并持有 `*sql.DB`，供每局复用，退出时在游戏循环结束后关闭。Loop 在原保存点通过一个可选回调同步写入，限定一次写入的时间，数据库细节留在 mineserver；保留进程内记录，写入失败记录 game_id 和错误，不中断下一局，不自动重试。同步写入会让下一局等本次写库结束或超时，首版不引入后台队列。
- 用户已在 `cmd/mineserver/match_store.go` 完成 `saveMatchResult`：同一事务先写对局、再逐人写参与者，全部成功才提交。定向构建通过；本机 MySQL 实测 1 局与 2 位玩家成功入库，重复对局编号返回 1062 且保留原数据，第二位玩家编号冲突时两表均为 0 行。诊断数据已清理，MySQL 已正常停止。见[验证记录](evidence/a3-match-write-20261119-02/verification.log)。
- 句柄交接已通过：用户修正 `probeMySQL` 成功返回的句柄，失败则关闭；`main` 持有并安排退出时关闭。定向构建通过，本机 MySQL 开启时用探测返回的句柄执行新的 `SELECT DATABASE()` 成功，关闭时返回空句柄和连接错误，核验后实例正常停止。[原始记录](evidence/a3-mysql-handle-20261119-01/verification.log)。
- 结算接线已由用户完成：`main` 提供带 3 秒超时的可选结算回调，`Loop` 在原内存保存点同步调用，失败记录对局编号和错误。静态审查发现两套循环争用通道，用户按定位删除了旧启动调用，正常停机只等唯一 `Loop` 关闭 `loopDone`，再由 `main` 退出关闭数据库。
- 本机验收：`mineserver` 定向构建与 `cmd/mineserver`、`internal/minesweeper` 原有测试通过。真实服务进程与两条本机 WebSocket 连接开 Coop / Easy，双人离开触发取消；SQL 查到同一编号的一条对局记录与两条玩家记录，`won=false`、`canceled=true`，两人 `left_early=true`。服务回到大厅，临时数据已清理，MySQL 已正常停止。见[构建与包测试](evidence/a3-loop-settlement-20261121-01/verification.log)和[联机观察](evidence/a3-loop-settlement-20261121-02/verification.log)。本轮未联机核验胜负、写库失败后继续、正常停机等待或重启查询；下一步设计并实现详情和有限分页，最后合并必要综合验收。
- 按用户要求，扫雷服务入口、游戏循环和事务保存源码中的说明性注释，以及上述源码的日志前缀和提示已中文化。内存保存日志标明是内存记录。定向构建、两个相关包测试及无数据库配置的启动冒烟通过；共享 Hub 日志与外部错误详情保留，历史验证日志未改。见[中文输出实测](evidence/a3-chinese-messages-20261121-01/verification.log)。
- 用户已按中文 TODO 完成 `cmd/mineserver/match_store.go` 的 `loadMatchResult`，先按主键查对局，再按同一编号查玩家并按 `player_id` 升序；AI 将 SQL 中写错的 `player_name` 改为表中已有列 `name`。定向构建及本机 MySQL 观察通过：完整详情与输入一致（含中文名字），玩家按编号排序，未知编号返回 `sql.ErrNoRows`。诊断数据已清理，MySQL 已正常停止。[原始记录](evidence/a3-match-detail-20261121-01/verification.log)。分页、对外查询入口及重启后查询仍待实现和验证。
- 用户按中文 TODO 完成 `listMatchResults`：仅查对局摘要，按唯一的 `game_id` 升序并使用 `LIMIT`/`OFFSET`，页大小 1 到 50、偏移量非负。随机编号排序不代表时间顺序，分页期间插入新数据可能改变页边界。AI 审查并用本机 MySQL 验证：相邻两页按序无重叠，50 条上限生效、非法参数被拒绝，合法空页返回空切片和无错误。定向构建通过，测试数据已清理，MySQL 已正常停止。[原始记录](evidence/a3-match-pagination-20261121-01/verification.log)。下一步是接详情与列表查询入口，再做重启后查询等综合验收。

- HTTP 查询入口从详情开始，复用存储层 `loadMatchResult`，计划使用 `GET /matches/{gameID}`。用户已按中文 TODO 完成路径取值、编号校验、查询及 JSON 响应，并将查无结果的 404 与其他查询错误的 500 分开。AI 修正响应写入失败时的日志格式，定向构建与包内既有测试通过。下一步补查询失败日志，并在 `main` 中仅在已连接数据库时注册详情路由。详情和分页接通后，验证重启后能通过 HTTP 查到已保存的结果，并执行全库测试与构建。
- 2026-11-23 详情接线审查：用户已保存 `db != nil` 和详情路由路径。详情路由误调用用于 `/ws` 的 `mineWSHandler`，`go build -mod=readonly -o <仓库外产物> ./cmd/mineserver` 退出 1，编译器报告该函数收到 `*sql.DB`，实际需要 Hub 和两个通道。已定位到详情这一行，待用户换成对应的处理函数后复验。
- 2026-11-23 详情接线复验：用户已将详情路由接到 `matchDetailHandler(db)`，定向构建与包内既有测试通过。自动执行策略拦截了通过 PowerShell 脚本启动进程并发 HTTP 请求的验证，实际 HTTP 响应待综合验收。AI 已准备 `matchListHandler` 的中文 TODO 和路由注册位置；首版 `GET /matches?limit=2&offset=0` 明确接收两个整数，缺少参数或超出范围时返回 400。用户下一步完成 URL 参数读取、调用已有 `listMatchResults` 和 JSON 响应，并一起补齐查询失败日志。
- 2026-11-23 用户确认今后先看真实结果和圈注，再回到代码与面试复盘。按此方式复测 A3 的 `listMatchResults`：本机 MySQL 中每次取两条，偏移 0 得到 01/02，偏移 2 得到 03/04。观察工具清理了 51 条测试数据，MySQL 正常停止。另运行 Go 小例子展示网址参数从字符串转成整数。当前 `matchListHandler` 仍是 TODO 骨架，先处理网址参数读取与校验，之后完成查询调用、响应和路由验收。[圈注图](evidence/a3-pagination-walkthrough-20261123-01/annotated.png)
- 2026-11-23 分页参数草稿审查：用户把独立 `copies` 示例保存进项目处理函数，`go test -mod=readonly ./cmd/mineserver -count=1 -v -timeout=30s` 在编译阶段退出 1，报 `request` 未定义，`count` 和 `err` 未使用。项目请求变量为 `r`，目标参数为 `limit`、`offset`。用户明确截图须来自实际测试结果窗口，AI 将本次原始输出在 VS Code 中打开并截图。[原始日志](evidence/a3-http-pagination-review-20261123-01/verification.log) · [现场截图](evidence/a3-http-pagination-review-20261123-01/native-test-output.jpg)
- 2026-11-23 暂停带学时复查：用户已保存 `limit`、`offset` 的解析和 `listMatchResults` 调用。`go test -mod=readonly ./cmd/mineserver -count=1 -v -timeout=30s` 在编译阶段退出 1，原因是第 56 行声明的 `page` 尚未用于 JSON 响应。待完成分页 JSON 响应、两个错误分支的 400/500 与提前返回、查询失败日志，以及注册 `GET /matches`。用户要求先检查并暂停带学，实际 HTTP 与重启综合验收待续。
- 2026-11-23 用户要求将当前进度提交并推送至个人 `origin/main`。推送前核对：用户已补好分页 400/500 和提前返回；分页 JSON 响应与路由注册仍待完成。定向测试因 `page` 尚未用于响应而在编译阶段退出 1。此批次保存暂停状态，A3 最终验收待续。[核对记录](evidence/a3-checkpoint-20261123-01/verification.log)
- 用户同意交付「大厅 → 对局战绩列表 → 单局详情」，并在指导下保存 `matchListHandler` 的分页 JSON 响应和 `main` 的列表路由。AI 补上路由中 `GET` 后遗漏的空格，复用玩家现有战绩详情页。A3 本机 HTTP 列表、分页和详情均已联调，浏览器实际打完一局能从列表点开，独立验证实例重启后读到同一局。全库 43 项测试与构建通过，在线试玩入口 `http://127.0.0.1:8093/` 保留。本机事务重复及回滚采用此前已验证的记录。[前端验证记录](evidence/a3-browser-history-20261123-01/verification.log) · [综合验收](evidence/a3-acceptance-20261123-01/verification.log)

## 开发反馈方式试行稿

1. 为每个有意义的功能节点约定一个玩家能完成的动作，以及一个能当场看到的结果。例如打完一局，点开战绩，找到这一局和自己的翻格数。
2. 提供一个能反复打开的试玩入口，并尽量使用刚保存的代码和真实数据。AI 负责启动、准备必要的测试条件以及操作第二个玩家。用户可以自己点、和 AI 一起玩，或先看 AI 操作一遍再接手。模拟数据要在页面或说明中标清来源。
3. 重复验证可以直接用脚本调用服务，快速准备操作或数据。达到功能节点时，用真实入口展示结果，用户可以接着点。页面已经在眼前就直接看页面，按需截图。用户随时能用一句话确认「我刚才做成了什么、下一步还能看到什么」。
4. 其他项目可以换成自己的可操作入口，例如 Unity 场景、浏览器页面或能发出实际请求的客户端。先在项目 A 试行，跨项目是否采用由用户之后决定。

## 项目 A 浏览器试玩（2026-11-23）

- 用户明确授权 AI 制作本机扫雷试玩客户端，由 AI 负责网页、连接和测试。首个验收点是两个独立浏览器标签页连接同一 A2 服务，开合作／简单局后翻格、插旗双向可见。网页在 `tools/playtest/`，启动时导出 `e3413e3` 到仓库外，只在副本加入网页入口和本机监听，A3 后端仍处于暂停快照。浏览器接收原协议的大厅、棋盘、光标及胜负消息。
- 首次本机联机已验证：两名浏览器测试玩家同局，第 5 行第 5 列翻开、第 1 行第 1 列插旗均同步，双方剩余旗数为 9；随后踩雷，双方回到大厅。使用 `http://127.0.0.1:8091/` 试玩。这次收尾保留了用户的试玩页以及已进入大厅的 GPT 玩家，等待用户加入。[结果与命令](evidence/a2-browser-playtest-20261123-01/verification.log) · [实际棋盘](evidence/a2-browser-playtest-20261123-01/duo-board.jpg)
- A3 第一段「打完一局，按同一编号查战绩」曾用 8092 详情预览通过。两名浏览器演示玩家打出一场取消局，MySQL 查询、详情接口和页面均显示演示甲翻 34 格、演示乙插 1 旗、两人中途离开。记录保留在本机学习库，也能通过目前 8093 的列表点开。早期 8092 仅在仓库外副本跳过分页草稿。[详情预览实测](evidence/a3-browser-detail-20261123-01/verification.log) · [最终联调实测](evidence/a3-acceptance-20261123-01/verification.log)

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
| R8 | A3 核对 Flags 字段时看到 Flag 插旗增加操作玩家的计数，撤旗按当前操作玩家递减，未按格子原 FlaggedBy 更新。 | 已于 2026-11-27 A4 验收中修复：用户在 `Flag` 的 `case StateFlagged:` 先保存 `cell.FlaggedBy` 再清零并按原插旗者扣减，补充 `TestCrossPlayerUnflagDecrementsOriginalOwner`，并在 `8093` 试玩页面实测修复前 `1`、修复后 `0`。 |

## Git 与个人仓库（2026-11-11）

- 用户已授权创建个人副本并使用本人 GitHub 仓库作为远程地址；AI 已通过登录为 LearnedYet 的 GitHub CLI 创建公开 Fork：[LearnedYet/multiplayer-terminal-games](https://github.com/LearnedYet/multiplayer-terminal-games)。GitHub API 确认 parent 为 siyad01/multiplayer-terminal-games。
- 本地仓库原已存在于 `C:\Dev\Go\A\.git`，未重新初始化。`origin` 指向个人 Fork，原作者地址保留为 `upstream`；默认推送远程为 origin。用户本轮授权本地提交，HEAD 已为 `d02f621`，未推送；详见 [提交与启动记录](evidence/a2-start-20261111-01/local-commit-and-phase-start.log)。
- 本次只创建上游副本、调整本地 Git 配置并记录事实，不修改 A1.2 业务代码或验收状态，也未进行实际 push。操作、首次参数错误及最终核对见 [Git 配置证据](evidence/git-personal-fork-20261111-01/setup.log)；这项配置不计为个人游戏功能实现。

## 下一步与未授权方向

A1.1/A1.2 历史成果保留；A2 已于 2026-11-14 完成停止，固定名单、连接身份隔离及必要综合验收已收尾，验收批次为 43 项全库测试及构建通过。2026-11-15 接续准备时只读核对：HEAD 为回执文档提交 702f469，工作区干净；本轮没有重跑测试或联网重查远程。此处原有“A2 尚未完整交付”的旧交接待办已按已验收事实修正。

A3 已完成本机工程验收并推送个人 `origin/main`（`36e97a5`）；A4 已验收并推送（`fc8c71c` / `94951ca`）。A5 双房间隔离已于 2026-11-28 完成本机验收（待用户确认是否提交推送至个人 `origin/main`），至此项目 A 首版达到路线图收尾终点，下一步按总路线可转入项目 B。

核心业务代码与核心规则小断言继续由用户按思路提示与 TODO 骨架填写，AI 负责调查、审查、前端/网络集成测试外壳与本地文档；后续阶段不自动提交、推送或部署。
