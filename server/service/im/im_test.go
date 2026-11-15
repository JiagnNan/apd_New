package im

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestAccountSyncAndUserSigUseSameUserAndSeparateCredentials(t *testing.T) {
	svc, err := New(Config{Enabled: true, SDKAppID: 123, AdminUser: "administrator", Region: "sg"}, "test-only-secret")
	if err != nil {
		t.Fatal(err)
	}
	var imported string
	var adminSignature string
	svc.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "adminapisgp.im.qcloud.com" || req.URL.Path != accountImportPath || req.Method != http.MethodPost {
			t.Fatal("unexpected Tencent account import endpoint")
		}
		if req.URL.Query().Get("sdkappid") != "123" || req.URL.Query().Get("identifier") != "administrator" {
			t.Fatal("expected account import to use administrator identity")
		}
		adminSignature = req.URL.Query().Get("usersig")
		var account struct{ UserID string }
		if err := json.NewDecoder(req.Body).Decode(&account); err != nil {
			t.Fatal(err)
		}
		imported = account.UserID
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ActionStatus":"OK","ErrorCode":0}`)),
		}, nil
	})

	if err := svc.SyncAccount(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	svc.client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("issuing UserSig must not call Tencent")
		return nil, nil
	})
	info, err := svc.IssueUserSig(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if imported != "42" || info.UserID != imported || info.SDKAppID != 123 || info.UserSig == adminSignature {
		t.Fatal("issued credentials must belong to the imported business user")
	}
	encoded := strings.NewReplacer("*", "+", "-", "/", "_", "=").Replace(info.UserSig)
	compressed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	var payload struct {
		UserID   string `json:"TLS.identifier"`
		IssuedAt int64  `json:"TLS.time"`
		TTL      int64  `json:"TLS.expire"`
	}
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.UserID != imported || payload.IssuedAt+payload.TTL != info.ExpiresAt {
		t.Fatal("UserSig payload does not match the issued identity and expiration")
	}
}

func TestIssueUserSigSucceedsAfterAccountImportFails(t *testing.T) {
	svc, err := New(Config{Enabled: true, SDKAppID: 123, AdminUser: "administrator"}, "test-only-secret")
	if err != nil {
		t.Fatal(err)
	}
	svc.client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"ActionStatus":"FAIL","ErrorCode":70169}`)),
		}, nil
	})
	if err := svc.SyncAccount(context.Background(), 42); err == nil {
		t.Fatal("account import failure must be reported to the background task")
	}
	if credentials, err := svc.IssueUserSig(context.Background(), 42); err != nil || credentials == nil {
		t.Fatalf("account import failure must not block local signing: %v", err)
	}
}

func TestIssueUserSigRejectsInvalidAndAdministratorIDs(t *testing.T) {
	svc, err := New(Config{Enabled: true, SDKAppID: 123, AdminUser: "42"}, "test-only-secret")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []int64{0, -1, 42} {
		if info, err := svc.IssueUserSig(context.Background(), id); err == nil || info != nil {
			t.Fatalf("issued credentials for invalid or administrator ID %d", id)
		}
	}
}

func TestAccountSyncDoesNotBlockLoginOrInheritRequestCancellation(t *testing.T) {
	svc, err := New(Config{Enabled: true, SDKAppID: 123, AdminUser: "administrator"}, "test-only-secret")
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan error, 1)
	release := make(chan struct{})
	defer close(release)
	svc.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		started <- req.Context().Err()
		select {
		case <-release:
			return &http.Response{StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{"ActionStatus":"OK","ErrorCode":0}`))}, nil
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	returned := make(chan struct{})
	go func() {
		svc.SyncAccountAsync(ctx, 42)
		close(returned)
	}()
	select {
	case <-returned:
	case <-time.After(time.Second):
		t.Fatal("account synchronization blocked the login caller")
	}
	select {
	case err := <-started:
		if err != nil {
			t.Fatalf("background sync inherited request cancellation: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("background synchronization did not start")
	}
}
