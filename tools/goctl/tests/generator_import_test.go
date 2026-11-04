package generator

import (
	"path/filepath"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
)

func TestLocalProtobufImport(t *testing.T) {
	root := t.TempDir()
	moduleRoot := filepath.Join(root, "server")
	for _, tc := range []struct{ name, declared, directory, want string }{
		{"cross module auth", "example.com/protocol/auth/authpb", filepath.Join(root, "protocol", "auth", "authpb"), "example.com/protocol/auth/authpb"},
		{"cross module user", "example.com/protocol/user/userpb", filepath.Join(root, "protocol", "user", "userpb"), "example.com/protocol/user/userpb"},
		{"same module", "example.com/server/rpc/auth/pb", filepath.Join(moduleRoot, "rpc", "auth", "pb"), "example.com/server/rpc/auth/pb"},
		{"relative compatibility", "./pb", filepath.Join(moduleRoot, "rpc", "auth", "pb"), "example.com/server/rpc/auth/pb"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := &defaultDirContext{inner: make(map[string]Dir), ctx: &ctx.ProjectContext{Dir: moduleRoot, Path: "example.com/server"}, pbImportPath: tc.declared}
			d.SetPbDir(tc.directory, tc.directory)
			if d.GetPb().Package != tc.want || d.GetProtoGo().Package != tc.want {
				t.Fatalf("got %q and %q, want %q", d.GetPb().Package, d.GetProtoGo().Package, tc.want)
			}
		})
	}
}
