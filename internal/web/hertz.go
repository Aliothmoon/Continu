package web

import (
	"context"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/Aliothmoon/Continu/internal/web/router"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Start() {
	h := server.Default(server.WithDisablePrintRoute(true))

	h.Use(handler.GlobalErrHandler())

	router.Register(h)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.Spin()
}
