// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package logic

import (
	"context"

	"test-go/go-zero/internal/svc"
	"test-go/go-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.ProfileResp, err error) {
	// todo: add your logic here and delete this line

	return
}
