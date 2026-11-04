param(
    [ValidatePattern('^[a-zA-Z][a-zA-Z0-9_-]*$')]
    [string]$Service = 'auth'
)

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent $PSScriptRoot
$protoPath = Join-Path $projectRoot "goprobuf/$Service/$Service.proto"
$rpcPath = Join-Path $projectRoot "server/rpc/$Service"

$protoText = [System.IO.File]::ReadAllText($protoPath)
$packageMatch = [regex]::Match($protoText, 'option\s+go_package\s*=\s*"([^"]+)"')
if (-not $packageMatch.Success) {
    throw "Missing go_package in $protoPath"
}
$goPackage = $packageMatch.Groups[1].Value
$packageParts = $goPackage.Split(';')
$importPath = $packageParts[0]
$packageName = ($importPath -split '/')[-1]
if ($packageParts.Length -gt 1) {
    $packageName = $packageParts[1]
}
# Older goctl versions sanitize the entire final go_package segment,
# including the semicolon alias, into e.g. authpb_authpb.
$legacyQualifier = (($goPackage -split '/')[-1] -replace '[^a-zA-Z0-9_]', '_')

$pbModuleText = [System.IO.File]::ReadAllText((Join-Path $projectRoot 'goprobuf/go.mod'))
$pbModule = [regex]::Match($pbModuleText, '(?m)^module\s+(\S+)').Groups[1].Value
$serverModuleText = [System.IO.File]::ReadAllText((Join-Path $projectRoot 'server/go.mod'))
$serverModule = [regex]::Match($serverModuleText, '(?m)^module\s+(\S+)').Groups[1].Value
if (-not $pbModule -or -not $serverModule -or -not $importPath.StartsWith($pbModule + '/')) {
    throw 'go_package must refer to a package inside the goprobuf module.'
}
$pbRelativePath = $importPath.Substring($pbModule.Length + 1)
$pbDirectory = (Resolve-Path -LiteralPath (Join-Path $projectRoot "goprobuf/$pbRelativePath")).Path
$badImport = $serverModule + '/' + $pbDirectory.Replace('\', '/')

$changed = 0
Get-ChildItem -LiteralPath $rpcPath -Recurse -Filter '*.go' | ForEach-Object {
    $original = [System.IO.File]::ReadAllText($_.FullName)
    $updated = $original.Replace('"' + $badImport + '"', '"' + $importPath + '"')
    if ($updated.Contains('"' + $importPath + '"') -and $legacyQualifier -ne $packageName) {
        $updated = $updated.Replace($legacyQualifier + '.', $packageName + '.')
    }
    if ($updated -cne $original) {
        [System.IO.File]::WriteAllText($_.FullName, $updated, [System.Text.UTF8Encoding]::new($false))
        $changed++
    }
}
Write-Output "Fixed $changed file(s) for $Service."
