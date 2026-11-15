package auth

import (
	"net/http"

	logic "github.com/JiagnNan/apd_New/api-gateway/internal/logic/auth"
	"github.com/JiagnNan/apd_New/api-gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RefreshImUserSigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		values := r.Header.Values("Authorization")
		if len(values) != 1 {
			writeIMUnauthorized(w)
			return
		}
		resp, err := logic.NewRefreshImUserSigLogic(r.Context(), svcCtx).RefreshImUserSig(values[0])
		if err != nil {
			if status.Code(err) == codes.Unauthenticated {
				writeIMUnauthorized(w)
				return
			}
			http.Error(w, "IM is unavailable", http.StatusServiceUnavailable)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func writeIMUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	http.Error(w, "valid access token required", http.StatusUnauthorized)
}
