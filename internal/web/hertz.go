package web

import (
	"context"
	_ "embed"
	"fmt"
	_ "github.com/Aliothmoon/Continu/internal/banner"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/Aliothmoon/Continu/internal/web/router"
	"github.com/cloudwego/hertz/pkg/app"
	auth "github.com/cloudwego/hertz/pkg/app/middlewares/server/basic_auth"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"os"
	"strconv"
	"strings"
)

var (
	//go:embed static/index.html
	index []byte
	//go:embed static/assets/index-72119a25.css
	css []byte
	c   = "/assets/index-72119a25.css"
	//go:embed static/assets/index-dde9c3ca.js
	js []byte
	j  = "/assets//index-dde9c3ca.js"
	//go:embed static/vite.svg
	svg []byte
	s   = "/vite.svg"
)

func Start() {

	h := server.Default(
		server.WithDisablePrintRoute(true),
		server.WithHostPorts(fmt.Sprintf(":%v", LoadPort())),
		server.WithExitWaitTime(0),
	)

	h.Use(handler.GlobalErrHandler())

	// Static Embed
	LoadFs(h)

	router.Register(h)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.Spin()
}

func Auth(h app.HandlerFunc) []app.HandlerFunc {
	env, ok := os.LookupEnv("CI_USER")
	ac := make(map[string]string)
	if ok {
		res := strings.Split(env, "#")
		if len(res) == 2 {
			ac[res[0]] = res[1]
		}
	}
	return []app.HandlerFunc{
		auth.BasicAuth(ac),
		h,
	}
}

func LoadFs(h *server.Hertz) {
	h.Any("/", Auth(func(c context.Context, ctx *app.RequestContext) {
		ctx.Redirect(consts.StatusTemporaryRedirect, []byte("/index.html"))
	})...)
	h.Any("/index.html", Auth(func(c context.Context, ctx *app.RequestContext) {
		ctx.Data(consts.StatusOK, consts.MIMETextHtml, index)
	})...)
	h.Any(c, Auth(func(c context.Context, ctx *app.RequestContext) {
		ctx.Data(consts.StatusOK, consts.MIMETextCss, css)
	})...)
	h.Any(j, Auth(func(c context.Context, ctx *app.RequestContext) {
		ctx.Data(consts.StatusOK, consts.MIMETextJavascript, js)
	})...)
	h.Any(s, Auth(func(c context.Context, ctx *app.RequestContext) {
		ctx.Data(consts.StatusOK, consts.MIMEImageSVG, svg)
	})...)
}

func LoadPort() int {
	env, ok := os.LookupEnv("CI_PORT")
	port := 6400
	if ok {
		var err error
		port, err = strconv.Atoi(env)
		if err != nil {
			hlog.Info("Load [CI_PORT] port fail , use default :6400")
		}
	}
	return port
}
