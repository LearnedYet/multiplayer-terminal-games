param([string]$Sql)
$ErrorActionPreference = 'Stop'
$mysqlArguments = @("--defaults-file=$PSScriptRoot\private\learning-client.cnf", '--database=go_a_learning')
if ($PSBoundParameters.ContainsKey('Sql')) { $mysqlArguments += @('--batch', "--execute=$Sql") }
& (Join-Path $PSScriptRoot 'mysql-8.4.13-winx64\bin\mysql.exe') @mysqlArguments
if ($LASTEXITCODE -ne 0) { throw "MySQL client failed: $LASTEXITCODE" }
