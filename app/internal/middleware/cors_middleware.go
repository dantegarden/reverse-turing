package middleware

import "net/http"

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//解决跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// Passthrough to next handler if need
		next(w, r)
	}
}
