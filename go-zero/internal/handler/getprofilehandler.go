// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"test-go/go-zero/internal/logic"
	"test-go/go-zero/internal/svc"
)

func GetProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetProfile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
