package meta

import (
	"fmt"
	"net/http"

	"reverse-turing/app/internal/logic/meta"
	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

// 与AI角色对话
func GameCharacterTalkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameCharacterTalkReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"code": 41400,
				"msg":  fmt.Sprintf("%s[%s]", "参数绑定/校验错误", err.Error()),
			})
			return
		}
		w.Header().Set(`Content-Type`, `text/event-stream;charset=utf-8`)
		l := meta.NewGameCharacterTalkLogic(r.Context(), svcCtx)
		_, err := l.GameCharacterTalk(&req, w)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
