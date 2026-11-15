package im

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	Enabled   bool   `json:",optional"`
	SDKAppID  int64  `json:",optional"`
	AdminUser string `json:",optional"`
	Region    string `json:",optional"`
}

type Service struct {
	cfg       Config
	secret    string
	client    *http.Client
	syncSlots chan struct{}
	importURL string
}

type LoginInfo struct {
	SDKAppID  int64
	UserID    string
	UserSig   string
	ExpiresAt int64
}

const (
	accountImportPath = "/v4/im_open_login_svc/account_import"
	userSigTTL        = int64(3600)
)

var regionHosts = map[string]string{
	"cn": "console.tim.qq.com",
	"sg": "adminapisgp.im.qcloud.com",
	"kr": "adminapikr.im.qcloud.com",
	"jp": "adminapijpn.im.qcloud.com",
	"de": "adminapiger.im.qcloud.com",
	"us": "adminapiusa.im.qcloud.com",
	"id": "adminapiidn.im.qcloud.com",
	"sa": "adminapiksa.im.qcloud.com",
}

func New(cfg Config, secret string) (*Service, error) {
	if !cfg.Enabled {
		return &Service{}, nil
	}
	if cfg.Region == "" {
		cfg.Region = "cn"
	}
	host, ok := regionHosts[cfg.Region]
	if !ok {
		return nil, errors.New("unsupported Tencent IM region")
	}
	if cfg.SDKAppID <= 0 || strings.TrimSpace(secret) == "" || strings.TrimSpace(cfg.AdminUser) == "" {
		return nil, errors.New("Tencent IM requires SDKAppID, AdminUser and APD_IM_SECRET_KEY")
	}
	return &Service{
		cfg: cfg, secret: secret,
		syncSlots: make(chan struct{}, 8),
		importURL: "https://" + host + accountImportPath,
		client: &http.Client{
			Timeout:       2 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}, nil
}

func (s *Service) Enabled() bool {
	return s != nil && s.cfg.Enabled
}

// SyncAccountAsync is called after a successful login or committed registration.
// A later login can retry an unsuccessful import with the same user ID.
func (s *Service) SyncAccountAsync(ctx context.Context, userID int64) {
	if !s.Enabled() {
		return
	}
	select {
	case s.syncSlots <- struct{}{}:
	default:
		logx.WithContext(ctx).Errorf("Tencent IM account sync busy for user %d", userID)
		return
	}
	go func() {
		defer func() { <-s.syncSlots }()
		syncCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		if err := s.SyncAccount(syncCtx, userID); err != nil {
			logx.WithContext(syncCtx).Errorf("sync Tencent IM account for user %d: %v", userID, err)
		}
	}()
}

// SyncAccount imports the committed local user ID into Tencent IM. It can be
// called again on later logins with the same ID.
func (s *Service) SyncAccount(ctx context.Context, userID int64) error {
	if !s.Enabled() {
		return nil
	}
	if userID <= 0 {
		return errors.New("invalid IM user ID")
	}
	id := strconv.FormatInt(userID, 10)
	if id == s.cfg.AdminUser {
		return errors.New("IM user ID collides with administrator")
	}
	adminSig, _, err := s.sign(s.cfg.AdminUser)
	if err != nil {
		return err
	}
	var nonce [4]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return err
	}
	query := url.Values{
		"sdkappid":    {strconv.FormatInt(s.cfg.SDKAppID, 10)},
		"identifier":  {s.cfg.AdminUser},
		"usersig":     {adminSig},
		"random":      {strconv.FormatUint(uint64(binary.BigEndian.Uint32(nonce[:])), 10)},
		"contenttype": {"json"},
	}
	body, err := json.Marshal(struct {
		UserID string `json:"UserID"`
	}{UserID: id})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.importURL+"?"+query.Encode(), bytes.NewReader(body))
	if err != nil {
		return errors.New("build Tencent IM account import request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		// net/http transport errors can contain the URL and administrator UserSig.
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("Tencent IM account import request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Tencent IM account import HTTP status %d", resp.StatusCode)
	}
	var result struct {
		ActionStatus string
		ErrorCode    int
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&result); err != nil {
		return errors.New("invalid Tencent IM account import response")
	}
	if result.ActionStatus != "OK" || result.ErrorCode != 0 {
		return fmt.Errorf("Tencent IM account import rejected with code %d", result.ErrorCode)
	}
	return nil
}

// IssueUserSig signs locally. Account synchronization runs asynchronously after
// login and must not prevent an existing IM user from refreshing credentials.
func (s *Service) IssueUserSig(ctx context.Context, userID int64) (*LoginInfo, error) {
	if !s.Enabled() {
		return nil, errors.New("Tencent IM is disabled")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if userID <= 0 {
		return nil, errors.New("invalid IM user ID")
	}
	id := strconv.FormatInt(userID, 10)
	if id == s.cfg.AdminUser {
		return nil, errors.New("IM user ID collides with administrator")
	}
	sig, expiresAt, err := s.sign(id)
	if err != nil {
		return nil, err
	}
	return &LoginInfo{SDKAppID: s.cfg.SDKAppID, UserID: id, UserSig: sig, ExpiresAt: expiresAt}, nil
}

func (s *Service) sign(id string) (string, int64, error) {
	now := time.Now().Unix()
	mac := hmac.New(sha256.New, []byte(s.secret))
	fmt.Fprintf(mac, "TLS.identifier:%s\nTLS.sdkappid:%d\nTLS.time:%d\nTLS.expire:%d\n", id, s.cfg.SDKAppID, now, userSigTTL)
	raw, err := json.Marshal(struct {
		Version    string `json:"TLS.ver"`
		Identifier string `json:"TLS.identifier"`
		SDKAppID   int64  `json:"TLS.sdkappid"`
		Time       int64  `json:"TLS.time"`
		Expire     int64  `json:"TLS.expire"`
		Signature  string `json:"TLS.sig"`
	}{
		Version: "2.0", Identifier: id, SDKAppID: s.cfg.SDKAppID, Time: now,
		Expire: userSigTTL, Signature: base64.StdEncoding.EncodeToString(mac.Sum(nil)),
	})
	if err != nil {
		return "", 0, err
	}
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	if _, err := writer.Write(raw); err != nil {
		_ = writer.Close()
		return "", 0, err
	}
	if err := writer.Close(); err != nil {
		return "", 0, err
	}
	sig := strings.NewReplacer("+", "*", "/", "-", "=", "_").Replace(base64.StdEncoding.EncodeToString(buf.Bytes()))
	return sig, now + userSigTTL, nil
}
