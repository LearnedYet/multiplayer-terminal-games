$ErrorActionPreference = 'Stop'
$mysqlRoot = $PSScriptRoot
$mysqlBin = Join-Path $mysqlRoot 'mysql-8.4.13-winx64\bin'
if (Get-NetTCPConnection -State Listen -LocalPort 13306 -ErrorAction SilentlyContinue) {
    throw 'Port 13306 is already listening; do not start a second instance.'
}
$mysqlProcess = Start-Process -FilePath (Join-Path $mysqlBin 'mysqld.exe') `
    -ArgumentList "--defaults-file=$mysqlRoot\my.ini" -WindowStyle Hidden -PassThru `
    -RedirectStandardOutput (Join-Path $mysqlRoot 'logs\process.stdout.log') `
    -RedirectStandardError (Join-Path $mysqlRoot 'logs\process.stderr.log')
for ($attempt = 0; $attempt -lt 20; $attempt++) {
    Start-Sleep -Milliseconds 250
    & (Join-Path $mysqlBin 'mysqladmin.exe') "--defaults-file=$mysqlRoot\private\admin-client.cnf" --connect-timeout=1 ping 2>$null | Out-Null
    if ($LASTEXITCODE -eq 0) {
        $mysqlServerId = [int](Get-Content -LiteralPath (Join-Path $mysqlRoot 'mysqld.pid'))
        $mysqlServerId | Set-Content -LiteralPath (Join-Path $mysqlRoot 'started-process-id.txt') -Encoding ascii
        Write-Output "MySQL ready at 127.0.0.1:13306; server PID=$mysqlServerId"
        return
    }
}
throw 'MySQL did not become ready; inspect the local logs directory.'
