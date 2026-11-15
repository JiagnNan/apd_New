package mdata

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strings"
)

// mdata 在 context 和 gRPC metadata 中使用的统一键名。
const dataKey = "mdata"

// Mdata 保存请求链路中的业务上下文信息。
// 后续如果改为 protobuf，只需要替换序列化方式并保留 FromContext 调用入口。
type Mdata struct {
	UserId         int64  `json:"userId"`
	DeviceModel    string `json:"deviceModel"`
	DeviceId       string `json:"deviceId"`
	Platform       string `json:"platform,omitempty"`
	SystemVersion  string `json:"systemVersion"`
	AppVersion     string `json:"appVersion"`
	SystemLangCode string `json:"systemLangCode"`
	LangCode       string `json:"langCode"`
	ClientIp       string `json:"clientIp"`
	LangPack       string `json:"langPack"`
	CountryCode    int64  `json:"countryCode"`
	Jti            string `json:"jti"`
	TokenType      string `json:"tokenType"`
	TokenHash      string `json:"tokenHash"`
	TokenExpiresAt int64  `json:"tokenExpiresAt"`
	UserAgent      string `json:"userAgent"`
}

// 将 Mdata 写入当前进程内的 context。
// 将 Mdata 结构体保存到 context.Context 中，方便后续函数通过 Context 获取当前请求的用户信息、设备信息等数据
func WithMdata(ctx context.Context, md *Mdata) context.Context {
	return context.WithValue(ctx, dataKey, md)
}

// 将 Mdata 同时写入 context 和 gRPC outgoing metadata，方便 API 调 RPC 时透传。
func ToContext(ctx context.Context, md *Mdata) context.Context {
	if md == nil {
		md = &Mdata{}
	}
	ctx = WithMdata(ctx, md)
	b, err := json.Marshal(md)
	if err != nil {
		return ctx
	}
	//获取当前 Context 中已有的出站 Metadata。
	out, ok := metadata.FromOutgoingContext(ctx)

	if !ok {
		out = metadata.MD{}
	} else {
		out = out.Copy()
	}
	out.Set(dataKey, base64.StdEncoding.EncodeToString(b))
	//创建一个带有指定出站 Metadata 的新 Context。
	return metadata.NewOutgoingContext(ctx, out)

}

// 从 context 中读取 Mdata，读取失败时返回空对象，避免调用方频繁判空。
func FromContext(ctx context.Context) *Mdata {
	md, err := GetMdata(ctx)
	if err != nil {
		return &Mdata{}
	}
	return md
}

// 从 context 或 metadata 中读取 Mdata，读取不到时返回错误。
func GetMdata(ctx context.Context) (*Mdata, error) {
	if v, ok := ctx.Value(dataKey).(*Mdata); ok && v != nil {
		return v, nil
	}
	if md, ok := fromMdata(ctx); ok {
		return md, nil
	}
	return nil, fmt.Errorf("mdata not found")

}

// 获取当前登录用户 ID。
func GetUserId(ctx context.Context) int64 {
	return FromContext(ctx).UserId
}

// 获取客户端 IP。
func GetClientIp(ctx context.Context) string {
	return FromContext(ctx).ClientIp
}

// 从 gRPC metadata 中解码 Mdata。
func fromMdata(ctx context.Context) (*Mdata, bool) {
	//从 context.Context 中读取客户端发送过来的 Metadata（元数据）
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = metadata.FromOutgoingContext(ctx)
	}
	if !ok {
		return &Mdata{}, false
	}
	values := md.Get(dataKey)
	if len(values) == 0 {
		return &Mdata{}, false
	}

	b, err := base64.StdEncoding.DecodeString(values[0])
	if err != nil {
		return &Mdata{}, false
	}
	var data Mdata

	if err = json.Unmarshal(b, &data); err != nil {
		return &Mdata{}, false
	}
	return &data, true
}

// 判断端类型
func fixLangPack(md *Mdata) {
	if md.LangPack == "" {
		deviceModel := strings.ToLower(md.DeviceModel)
		systemVersion := strings.ToLower(md.SystemVersion)
		switch {
		case strings.HasPrefix(deviceModel, "macbook") ||
			strings.HasPrefix(systemVersion, "macos") ||
			strings.HasPrefix(systemVersion, "windows"):
			md.LangPack = "tdesktop"
		case strings.HasPrefix(deviceModel, "iphone"):
			md.LangPack = "ios"
		default:
			md.LangPack = "android" //默认android
		}
	}
}
