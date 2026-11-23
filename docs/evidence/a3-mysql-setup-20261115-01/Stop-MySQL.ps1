$ErrorActionPreference = 'Stop'
$mysqlServerId = [int](Get-Content -LiteralPath (Join-Path $PSScriptRoot 'mysqld.pid'))
$mysqlServer = Get-Process -Id $mysqlServerId -ErrorAction Stop
& (Join-Path $PSScriptRoot 'mysql-8.4.13-winx64\bin\mysqladmin.exe') `
    "--defaults-file=$PSScriptRoot\private\admin-client.cnf" --connect-timeout=5 shutdown
if ($LASTEXITCODE -ne 0) { throw "MySQL shutdown failed: $LASTEXITCODE" }
if (-not $mysqlServer.WaitForExit(15000)) { throw 'MySQL has not exited after shutdown; inspect the local logs.' }
Write-Output "MySQL server PID=$mysqlServerId exited normally."
