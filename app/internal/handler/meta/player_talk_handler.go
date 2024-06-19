package meta

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"reverse-turing/app/internal/logic/meta"
	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"
)

// 用户对话
func PlayerTalkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GamePlayerTalkReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := meta.NewPlayerTalkLogic(r.Context(), svcCtx)
		resp, err := l.PlayerTalk(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
