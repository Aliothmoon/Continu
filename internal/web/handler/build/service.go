package build

import (
	"context"
	"errors"
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gen"
	"strconv"
	"sync"
)

const (
	RID = "RID"
)

var (
	DProject   = query.Project
	DRecord    = query.BuildRecord
	ProcessMap sync.Map
)

func GetBuildHistoryList(c context.Context, ctx *app.RequestContext) {
	pid, err := strconv.Atoi(ctx.Param(biz.PID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	records, err := DRecord.Where(DRecord.Pid.Eq(int32(pid))).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Code: 0,
		Data: records,
	})
}

func AddBuildTask(c context.Context, ctx *app.RequestContext) {
	pid, err := strconv.Atoi(ctx.Param(biz.PID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	cond := []gen.Condition{DProject.ID.Eq(int32(pid)), DProject.Status.Eq(biz.ProjectIdle)}
	ps, err := DProject.Where(cond...).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if len(ps) == 0 {
		handler.LaunchError(ctx, errors.New("Record Not Found "))
		return
	}

	result, err := DProject.Where(cond...).Update(DProject.Status, biz.ProjectPending)
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if result.RowsAffected != 1 {
		handler.LaunchError(ctx, errors.New("Optimistic lock takes effect "))
		return
	}
	project := ps[0]
	var status int32 = biz.BuildPending
	record := model.BuildRecord{
		Pid:        &project.ID,
		Status:     &status,
		Branch:     project.Branch,
		Script:     project.Script,
		WorkDir:    project.WorkDir,
		ProjectURL: project.ProjectURL,
	}

	err = DRecord.Create(&record)
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}

	publishTask(record.ID, project)

}

func CancelBuildTask(c context.Context, ctx *app.RequestContext) {

}
