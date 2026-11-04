$ErrorActionPreference = 'Stop'
$patchedBinary = Join-Path $PSScriptRoot 'bin/goctl.exe'
if (-not (Test-Path -LiteralPath $patchedBinary)) { throw 'Run tools/goctl/build.ps1 first.' }
$target = (Get-Command goctl.exe -CommandType Application).Source
if ([System.IO.Path]::GetFullPath($target) -eq [System.IO.Path]::GetFullPath($patchedBinary)) {
    Write-Output 'The patched executable is already selected by PATH.'
    return
}
$backup = Join-Path (Split-Path -Parent $target) ('goctl.backup-' + [guid]::NewGuid().ToString('N') + '.exe')
Copy-Item -LiteralPath $target -Destination $backup
try {
    Copy-Item -LiteralPath $patchedBinary -Destination $target -Force
    $expectedHash = (Get-FileHash -LiteralPath $patchedBinary -Algorithm SHA256).Hash
    if ((Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash -ne $expectedHash) {
        throw 'Installed executable hash does not match.'
    }
} catch {
    Copy-Item -LiteralPath $backup -Destination $target -Force
    throw
}
Write-Output "Installed: $target"
Write-Output "Backup: $backup"
