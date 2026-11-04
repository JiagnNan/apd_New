package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalGoPackageAlias(t *testing.T) {
	for _, tc := range []struct{ option, path, name string }{
		{"example.com/protocol/auth/authpb;authpb", "example.com/protocol/auth/authpb", "authpb"},
		{"example.com/protocol/auth/authpb;custompb", "example.com/protocol/auth/authpb", "custompb"},
		{"example.com/protocol/user/userpb", "example.com/protocol/user/userpb", "userpb"},
		{"./pb", "./pb", "pb"},
	} {
		t.Run(tc.option, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "auth.proto")
			text := `syntax = "proto3"; package auth; option go_package = "` + tc.option + `";
message Req {} message Resp {} service Auth { rpc Login(Req) returns (Resp); }`
			if err := os.WriteFile(file, []byte(text), 0600); err != nil { t.Fatal(err) }
			got, err := NewDefaultProtoParser().Parse(file)
			if err != nil { t.Fatal(err) }
			if got.GoPackage != tc.path || got.PbPackage != tc.name {
				t.Fatalf("got %q;%q, want %q;%q", got.GoPackage, got.PbPackage, tc.path, tc.name)
			}
		})
	}
}
