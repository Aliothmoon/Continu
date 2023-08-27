package hooks

import (
	"context"
	"github.com/Aliothmoon/Continu/internal/logger"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/Aliothmoon/Continu/internal/web/handler/build"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

func ProcessWebHooks(c context.Context, ctx *app.RequestContext) {
	pid, err := strconv.Atoi(ctx.Param(biz.PID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if build.InternalProcessTask(pid, ctx) {
		return
	}
	logger.Infof("Received WebHooks 2 Project %v", pid)

	ctx.Status(consts.StatusOK)
}
