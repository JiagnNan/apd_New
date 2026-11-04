$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$binary = Join-Path $PSScriptRoot 'bin/goctl.exe'
$verifyRoot = Join-Path $PSScriptRoot ('verify/' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path (Join-Path $verifyRoot 'goprobuf/auth'), (Join-Path $verifyRoot 'goprobuf/user'), (Join-Path $verifyRoot 'server') -Force | Out-Null
Copy-Item (Join-Path $projectRoot 'goprobuf/auth/auth.proto') (Join-Path $verifyRoot 'goprobuf/auth/auth.proto')
Copy-Item (Join-Path $projectRoot 'goprobuf/go.mod') (Join-Path $verifyRoot 'goprobuf/go.mod')
# Reuse the project's dependency versions in an isolated fixture module.
$dependencies = [System.IO.File]::ReadAllText((Join-Path $projectRoot 'api-gateway/go.mod'))
$dependencies = [regex]::Replace($dependencies, '(?m)^module\s+\S+', 'module github.com/JiagnNan/apd_New/server')
[System.IO.File]::WriteAllText((Join-Path $verifyRoot 'server/go.mod'), $dependencies)
$workVersion = [regex]::Match([System.IO.File]::ReadAllText((Join-Path $projectRoot 'go.work')), '(?m)^go\s+(\S+)').Groups[1].Value
[System.IO.File]::WriteAllText((Join-Path $verifyRoot 'go.work'), "go $workVersion`nuse (`n ./server`n ./goprobuf`n)`n")
$userProto = @'
syntax = "proto3";
package user;
option go_package = "github.com/JiagnNan/apd_New/goprobuf/user/userpb;custompb";
message ProfileReq {} message ProfileResp {}
service User { rpc Profile(ProfileReq) returns (ProfileResp); }
'@
[System.IO.File]::WriteAllText((Join-Path $verifyRoot 'goprobuf/user/user.proto'), $userProto)
$oldGoWork = $env:GOWORK
$env:GOWORK = Join-Path $verifyRoot 'go.work'
Push-Location $verifyRoot
try {
    foreach ($service in @('auth', 'user')) {
        & $binary rpc protoc "./goprobuf/$service/$service.proto" `
            -I ./goprobuf `
            --go_out=./goprobuf `
            --go_opt=module=github.com/JiagnNan/apd_New/goprobuf `
            --go-grpc_out=./goprobuf `
            --go-grpc_opt=module=github.com/JiagnNan/apd_New/goprobuf `
            "--zrpc_out=./server/rpc/$service" `
            --module=github.com/JiagnNan/apd_New/server `
            --client=false
        if ($LASTEXITCODE -ne 0) { throw "Generation failed: $service" }
        $badFiles = Get-ChildItem "server/rpc/$service" -Recurse -Filter '*.go' |
            Select-String -Pattern '[A-Za-z]:[/\\]|authpb_authpb|userpb_custompb'
        if ($badFiles) { throw "Invalid generated imports: $badFiles" }
    }
    & go test ./server/rpc/auth/... ./server/rpc/user/...
    if ($LASTEXITCODE -ne 0) { throw 'Generated services did not compile.' }
    Write-Output "Fresh generation and compilation passed: $verifyRoot"
} finally {
    Pop-Location
    $env:GOWORK = $oldGoWork
}
