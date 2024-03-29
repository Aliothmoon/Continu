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
	"os"
	"strconv"
	"sync"
	"time"
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
	spec := pid != -1

	var records []*model.BuildRecord
	if spec {
		records, err = DRecord.Where(DRecord.Pid.Eq(int32(pid))).Order(DRecord.CreatedAt.Desc()).Find()
		if err != nil {
			handler.LaunchError(ctx, err)
			return
		}
	} else {
		records, err = DRecord.Order(DRecord.CreatedAt.Desc()).Find()
		if err != nil {
			handler.LaunchError(ctx, err)
			return
		}
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
	if InternalProcessTask(pid, ctx) {
		return
	}

	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Msg: "Add Task Complete",
	})

}

func InternalProcessTask(pid int, ctx *app.RequestContext) bool {
	cond := []gen.Condition{DProject.ID.Eq(int32(pid)), DProject.Status.Eq(biz.ProjectIdle)}
	ps, err := DProject.Where(cond...).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return true
	}
	if len(ps) == 0 {
		handler.LaunchError(ctx, errors.New("Record Not Found "))
		return true
	}

	result, err := DProject.Where(cond...).Update(DProject.Status, biz.ProjectPending)
	if err != nil {
		handler.LaunchError(ctx, err)
		return true
	}
	if result.RowsAffected != 1 {
		handler.LaunchError(ctx, errors.New("Optimistic lock takes effect "))
		return true
	}
	project := ps[0]
	var status int32 = biz.BuildPending
	record := model.BuildRecord{
		Pid:        project.ID,
		Status:     status,
		Parameters: project.Parameters,
		Bin:        project.Bin,
		WorkDir:    project.WorkDir,
		CreatedAt:  time.Now().UnixMilli(),
	}

	err = DRecord.Create(&record)
	if err != nil {
		handler.LaunchError(ctx, err)
		return true
	}
	go PublishTask(&ConstructInfo{
		BuildID:     record.ID,
		ProjectInfo: project,
		Log:         NewLogWriteCloser(record.ID),
	})
	return false
}

func CancelBuildTask(c context.Context, ctx *app.RequestContext) {
	rid, err := strconv.Atoi(ctx.Param(biz.RID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	value, ok := ProcessMap.Load(int32(rid))
	var msg string
	if ok {
		process := value.(*os.Process)
		err := Kill(process)
		process.Release()
		if err != nil {
			msg += err.Error()
		} else {
			msg = "Kill Ok"
		}
	} else {
		msg = "Process Not Found "
	}
	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Msg: msg,
	})
}
