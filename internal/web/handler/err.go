package handler

import (
	"context"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func LaunchError(c *app.RequestContext, err error) {
	_ = c.Error(err)
}
func GlobalErrHandler() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Next(c)
		last := ctx.Errors.Last()
		if last != nil && last.Err != nil {
			err := last.Err
			ctx.JSON(consts.StatusServiceUnavailable, biz.JsonModel{
				Code: 500,
				Msg:  err.Error(),
				Data: last.Err,
			})
		}
	}
}
