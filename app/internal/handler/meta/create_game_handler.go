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

// 创建游戏
func CreateGameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"code": 41400,
				"msg":  fmt.Sprintf("%s[%s]", "参数绑定/校验错误", err.Error()),
			})
			return
		}

		l := meta.NewCreateGameLogic(r.Context(), svcCtx)
		resp, err := l.CreateGame(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
