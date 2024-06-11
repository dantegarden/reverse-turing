package {{.PkgName}}

import (
        "fmt"
	"net/http"

	{{.ImportPackages}}
        
        "github.com/zeromicro/go-zero/rest/httpx"
        xhttp "github.com/zeromicro/x/http"
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
                    httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
		        "code": 41400,
			"msg":  fmt.Sprintf("%s[%s]", "参数绑定/校验错误", err.Error()),
		    })
                    return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
		    xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
		    xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
