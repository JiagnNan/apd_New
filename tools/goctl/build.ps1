$ErrorActionPreference = 'Stop'
$sourceVersion = 'v1.10.0'
$oldGoWork = $env:GOWORK
$env:GOWORK = 'off'
Push-Location $PSScriptRoot
try {
    $moduleCache = (& go env GOMODCACHE).Trim()
    if ($LASTEXITCODE -ne 0) { throw 'Cannot locate the Go module cache.' }
    $source = Join-Path $moduleCache "github.com/zeromicro/go-zero/tools/goctl@$sourceVersion"
    if (-not (Test-Path -LiteralPath $source)) {
        & go mod download "github.com/zeromicro/go-zero/tools/goctl@$sourceVersion"
        if ($LASTEXITCODE -ne 0) { throw 'Cannot download the pinned goctl source.' }
    }

    # Copy instead of modifying the shared, read-only Go module cache.
    $buildPath = Join-Path $PSScriptRoot '.build'
    New-Item -ItemType Directory -Path $buildPath -Force | Out-Null
    Get-ChildItem -LiteralPath $source -Force | Copy-Item -Destination $buildPath -Recurse -Force
    Get-ChildItem -LiteralPath $buildPath -Recurse -File -Force | ForEach-Object { $_.IsReadOnly = $false }

    function Update-Source([string]$RelativePath, [string]$Before, [string]$After) {
        $target = Join-Path $buildPath $RelativePath
        $content = [System.IO.File]::ReadAllText($target).Replace("`r`n", "`n")
        if (-not $content.Contains($Before)) { throw "Source anchor not found: $RelativePath" }
        [System.IO.File]::WriteAllText($target, $content.Replace($Before, $After), [System.Text.UTF8Encoding]::new($false))
    }

    $before = @'
	ret.PbPackage = GoSanitized(filepath.Base(ret.GoPackage))
'@
    $after = @'
	importPath, packageName, hasAlias := strings.Cut(ret.GoPackage, ";")
	ret.GoPackage = importPath
	if !hasAlias {
		packageName = filepath.Base(importPath)
	}
	ret.PbPackage = GoSanitized(packageName)
'@
    Update-Source 'rpc/parser/parser.go' $before $after

    Update-Source 'rpc/generator/mkdir.go' 'ctx         *ctx.ProjectContext' "ctx         *ctx.ProjectContext`n`t`tpbImportPath string"
    Update-Source 'rpc/generator/mkdir.go' 'ctx:         ctx,' "ctx:         ctx,`n`t`tpbImportPath: proto.GoPackage,"
    $before = @'
func (d *defaultDirContext) SetPbDir(pbDir, grpcDir string) {
	d.inner[pb] = Dir{
		Filename: pbDir,
		Package:  filepath.ToSlash(filepath.Join(d.ctx.Path, strings.TrimPrefix(pbDir, d.ctx.Dir))),
		Base:     filepath.Base(pbDir),
	}

	d.inner[protoGo] = Dir{
		Filename: grpcDir,
		Package: filepath.ToSlash(filepath.Join(d.ctx.Path,
			strings.TrimPrefix(grpcDir, d.ctx.Dir))),
		Base: filepath.Base(grpcDir),
	}
}
'@
    $after = @'
func (d *defaultDirContext) SetPbDir(pbDir, grpcDir string) {
	pbPackage := filepath.ToSlash(filepath.Join(d.ctx.Path, strings.TrimPrefix(pbDir, d.ctx.Dir)))
	grpcPackage := filepath.ToSlash(filepath.Join(d.ctx.Path, strings.TrimPrefix(grpcDir, d.ctx.Dir)))
	// A full go_package is already the canonical Go import path. It may
	// belong to a different module from the generated RPC service.
	if p := d.pbImportPath; strings.Contains(p, "/") && !strings.HasPrefix(p, ".") && !filepath.IsAbs(p) {
		pbPackage, grpcPackage = p, p
	}
	d.inner[pb] = Dir{Filename: pbDir, Package: pbPackage, Base: filepath.Base(pbDir)}
	d.inner[protoGo] = Dir{Filename: grpcDir, Package: grpcPackage, Base: filepath.Base(grpcDir)}
}
'@
    Update-Source 'rpc/generator/mkdir.go' $before $after
    Copy-Item (Join-Path $PSScriptRoot 'tests/parser_import_test.go') (Join-Path $buildPath 'rpc/parser/local_import_test.go') -Force
    Copy-Item (Join-Path $PSScriptRoot 'tests/generator_import_test.go') (Join-Path $buildPath 'rpc/generator/local_import_test.go') -Force

    Push-Location $buildPath
    try {
        & gofmt -w rpc/parser/parser.go rpc/generator/mkdir.go rpc/parser/local_import_test.go rpc/generator/local_import_test.go
        if ($LASTEXITCODE -ne 0) { throw 'Formatting failed.' }
        & go test ./rpc/parser ./rpc/generator -run 'TestLocal|TestDefaultProtoParse' -count=1
        if ($LASTEXITCODE -ne 0) { throw 'goctl regression tests failed.' }
        $binaryDir = Join-Path $PSScriptRoot 'bin'
        New-Item -ItemType Directory -Path $binaryDir -Force | Out-Null
        & go build -o (Join-Path $binaryDir 'goctl.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'Building goctl failed.' }
        Write-Output "Patched goctl: $binaryDir\goctl.exe (source $sourceVersion)"
    } finally {
        Pop-Location
    }
} finally {
    Pop-Location
    $env:GOWORK = $oldGoWork
}
