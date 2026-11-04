# Project goctl fix

This builds a source-level fix for RPC generation across the `server` and
`goprobuf` modules. Generated files do not need an import-rewriting step.

The source is pinned to `github.com/zeromicro/go-zero/tools/goctl@v1.10.0`,
which matches the module recorded in the installed executable's build metadata.
That upstream source still reports `1.9.2` in `goctl --version`.

The build script copies the source into the ignored `.build` directory, then:

- Separates the import path and optional package alias in `go_package`.
- Uses full `go_package` import paths when generating RPC imports, including
  protobuf packages outside the server module.
- Retains the original behavior for relative `go_package` paths.
- Runs regression tests and builds `bin/goctl.exe`.

Run from the repository root:

```powershell
.\tools\goctl\build.ps1
.\tools\goctl\verify.ps1
.\tools\goctl\install.ps1
```

Verification generates fresh auth and user services in an isolated, ignored
workspace and compiles both, including a custom protobuf package alias.
Installation backs up the executable currently selected by PATH before replacing
it. After installation, the existing `goctl rpc protoc ...` command works without
`scripts/fix-rpc-imports.ps1`.

An upstream `go install .../goctl@...` will replace the custom executable. Rebuild,
verify, and reinstall this fix if that version still has the problem. This is a
local customization, not an upstream release or an upstream-submitted patch.
