param(
    [ValidateRange(1024, 65535)]
    [int]$Port = 8091,
    [ValidateSet('A2', 'A3Detail', 'A3')]
    [string]$Backend = 'A2'
)

$ErrorActionPreference = 'Stop'
$repo = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$revision = 'e3413e3b66deaac9113f4ab1e94cd77dcbe59b18'
if ($Backend -eq 'A3Detail') {
    $revision = '8ae341221b32e7f1b034cc38e11a59f2a89356fe'
}
if ($Backend -ne 'A2') {
    if (-not $env:MINESWEEPER_MYSQL_DSN) { throw 'A3 试玩需要本机学习库连接。' }
}
$artifactRoot = Join-Path $env:LOCALAPPDATA 'GoA\artifacts\playtest'
$session = Join-Path $artifactRoot ([guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $session | Out-Null

if ($Backend -eq 'A3') {
    # A3 联调使用当前已保存的源码，用户完成分页后可直接试玩。
    $commandDir = New-Item -ItemType Directory -Path (Join-Path $session 'cmd')
    $sourceFiles = 'go.mod', 'go.sum', 'internal' | ForEach-Object { Join-Path $repo $_ }
    Copy-Item -LiteralPath $sourceFiles -Destination $session -Recurse
    Copy-Item -LiteralPath (Join-Path $repo 'cmd/mineserver') -Destination $commandDir.FullName -Recurse
} else {
    $archive = Join-Path $session 'backend.tar'
    git -C $repo archive --format=tar -o $archive $revision go.mod go.sum cmd/mineserver internal
    if ($LASTEXITCODE -ne 0) { throw '无法导出试玩服务源码。' }
    tar.exe -xf $archive -C $session
    if ($LASTEXITCODE -ne 0) { throw '无法解开试玩服务源码。' }
    Remove-Item -LiteralPath $archive
}

if ($Backend -eq 'A3Detail') {
    # 仅在临时副本去掉尚未注册的分页草稿，保留已有保存和详情链路。
    $handlerPath = Join-Path $session 'cmd\mineserver\match_http.go'
    $handlers = [IO.File]::ReadAllText($handlerPath)
    $draft = $handlers.IndexOf('func matchListHandler(')
    if ($draft -lt 0) { throw '找不到 A3 分页草稿。' }
    $detailOnly = $handlers.Substring(0, $draft) -replace '(?m)^[ \t]*"strconv"[ \t]*\r?\n', ''
    [IO.File]::WriteAllText($handlerPath, $detailOnly)
}

# 只给临时副本增加本机网页入口。原游戏规则和仓库源码沿用原样。
$entry = Join-Path $session 'cmd\mineserver\main.go'
$source = [IO.File]::ReadAllText($entry)
$anchor = 'srv := &http.Server{Addr: ":8081", Handler: nil}'
if (-not $source.Contains($anchor)) { throw '找不到试玩服务地址。' }
$localEntry = @'
http.Handle("/", http.FileServer(http.Dir(os.Getenv("GOA_PLAYTEST_PUBLIC"))))
	srv := &http.Server{Addr: os.Getenv("GOA_PLAYTEST_ADDR"), Handler: nil}
'@
$source = $source.Replace($anchor, $localEntry)
$source = $source.Replace('log.Printf("minesweeper server listening on :8081")', 'log.Printf("试玩服务：http://%s", srv.Addr)')
$source = $source.Replace('log.Printf("扫雷服务正在监听 :8081")', 'log.Printf("试玩服务：http://%s", srv.Addr)')
$source = $source.Replace('CheckOrigin: func(r *http.Request) bool { return true }', 'CheckOrigin: func(r *http.Request) bool { return r.Header.Get("Origin") == "http://" + os.Getenv("GOA_PLAYTEST_ADDR") }')
[IO.File]::WriteAllText($entry, $source)

$env:GOWORK = 'off'
$env:GOA_PLAYTEST_ADDR = "127.0.0.1:$Port"
$env:GOA_PLAYTEST_PUBLIC = Join-Path $PSScriptRoot 'public'
$binary = Join-Path $session 'minesweeper-playtest.exe'
Push-Location -LiteralPath $session
try {
    go build -mod=readonly -o $binary ./cmd/mineserver
    if ($LASTEXITCODE -ne 0) { throw '试玩服务构建失败。' }
    & $binary
    exit $LASTEXITCODE
} finally {
    Pop-Location
}
