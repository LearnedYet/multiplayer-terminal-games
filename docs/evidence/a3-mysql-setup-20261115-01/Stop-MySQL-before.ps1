$ErrorActionPreference = 'Stop'
& (Join-Path $PSScriptRoot 'mysql-8.4.13-winx64\bin\mysqladmin.exe') `
    "--defaults-file=$PSScriptRoot\private\admin-client.cnf" --connect-timeout=5 shutdown
if ($LASTEXITCODE -ne 0) { throw "MySQL shutdown failed: $LASTEXITCODE" }
Write-Output 'MySQL shutdown request completed.'
