package logs

import (
	"context"
	"errors"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
)

var (
	DLog    = query.Log
	DRecord = query.BuildRecord
)

func GetBuildLogs(c context.Context, ctx *app.RequestContext) {
	rid, err := strconv.Atoi(ctx.Param(biz.RID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	ts := ctx.Query(biz.TimeStamp)
	if ts == "" {
		handler.LaunchError(ctx, errors.New("Ts is Nil "))
		return
	}

	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}

	ti := time.UnixMilli(t)
	logs, err := DLog.Where(DLog.BuildID.Eq(int32(rid)), DLog.CreatedAt.Gt(ti)).Limit(300).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	var status int
	err = DRecord.Where(DRecord.ID.Eq(int32(rid))).Select(DRecord.Status).Scan(&status)
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}

	ls := make([]*biz.Log, len(logs))
	for i := range ls {
		l := logs[i]
		ls[i] = &biz.Log{
			ID:        l.ID,
			BuildID:   l.BuildID,
			Content:   l.Content,
			CreatedAt: l.CreatedAt.UnixMilli(),
		}
	}

	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Code: status,
		Data: ls,
	})

}
