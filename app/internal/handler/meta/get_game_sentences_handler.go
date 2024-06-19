package meta

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"reverse-turing/app/internal/logic/meta"
	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"
)

// 获取游戏对话记录
func GetGameSentencesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GamePathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := meta.NewGetGameSentencesLogic(r.Context(), svcCtx)
		resp, err := l.GetGameSentences(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
