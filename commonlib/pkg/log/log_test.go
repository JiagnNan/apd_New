package log

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	SetConfig(Config{Level: DebugLevel, Encoding: "console"})
	ctx := WithString(context.Background(), "userId", "1001")
	Debugx(ctx, "debug %s", "ok")
	Infox(ctx, "info")
	Warnf("warn")
	Errorf("error")
}
