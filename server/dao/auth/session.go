package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JiagnNan/apd_New/commonlib/redisx"
	"golang.org/x/net/context"
	"time"
)

type SessionDao struct {
	redisClient *redisx.Client
}

func NewSessionDao(client *redisx.Client) *SessionDao {
	return &SessionDao{
		redisClient: client,
	}
}

type Session struct {
	SessionId       string `json:"sid"`
	AccessHash      string `json:"accessHash"`
	RefreshHash     string `json:"refreshHash"`
	AuthVersion     string `json:"authVersion"`
	RefreshExpireAt int64  `json:"refreshExpireAt"`
	DeviceId        string `json:"deviceId,optional"`
}

const (
	sessionKeyFormat = "auth:session:%d:%s"
)

func (d *SessionDao) GetSession(ctx context.Context, userID int64, platform string) (*Session, error) {
	raw, err := d.redisClient.Get(ctx, fmt.Sprintf(sessionKeyFormat, userID, platform))
	if err != nil {
		return nil, err
	}
	var session Session
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (d *SessionDao) SaveSession(ctx context.Context, userId int64, platform string, session Session) error {
	key := fmt.Sprintf(sessionKeyFormat, userId, platform)
	t := time.Until(time.Unix(session.RefreshExpireAt, 0))
	if t <= 0 {
		return errors.New("session expired")
	}
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return d.redisClient.Set(ctx, key, b, t)
}
