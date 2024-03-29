package project

import (
	"context"
	"errors"
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"github.com/Aliothmoon/Continu/internal/repo/query"
	"github.com/Aliothmoon/Continu/internal/web/biz"
	"github.com/Aliothmoon/Continu/internal/web/handler"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
)

var (
	DProject = query.Project
	DRecord  = query.BuildRecord
	DLog     = query.Log
)

func AddProject(c context.Context, ctx *app.RequestContext) {
	var p biz.Project
	if err := ctx.Bind(&p); err != nil {
		handler.LaunchError(ctx, err)
		return
	}

	if p.Name == "" {
		handler.LaunchError(ctx, errors.New("Project Not Null "))
		return
	}

	var isGit int32 = biz.GitProject
	if !p.IsGit {
		isGit = biz.NoneGitProject
	}
	err := DProject.Create(&model.Project{
		Name:       p.Name,
		Status:     p.Status,
		Branch:     &p.Branch,
		ProjectURL: &p.ProjectURL,
		WorkDir:    &p.WorkDir,
		PrivateKey: &p.PrivateKey,
		IsGit:      &isGit,
		Bin:        &p.Bin,
		Parameters: &p.Parameters,
		CreatedAt:  time.Now().UnixMilli(),
	})
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Msg: "Create Ok",
	})
}

func DelProject(c context.Context, ctx *app.RequestContext) {
	pid, err := strconv.Atoi(ctx.Param(biz.PID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	info, err := DProject.Where(DProject.ID.Eq(int32(pid))).Delete()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if info.RowsAffected == 0 {
		handler.LaunchError(ctx, errors.New("Delete Failed Can't Find ProjectInfo "))
		return
	}
	var ids []int32
	err = DRecord.Where(DRecord.Pid.Eq(int32(pid))).Select(DRecord.ID).Scan(&ids)
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	_, err = DRecord.Where(DRecord.Pid.Eq(int32(pid))).Delete()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	_, err = DLog.Where(DLog.BuildID.In(ids...)).Delete()

	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}

	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Msg: "Del Ok",
	})
}

func UpdateProject(c context.Context, ctx *app.RequestContext) {
	var p biz.Project
	if err := ctx.Bind(&p); err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if p.Name == "" {
		handler.LaunchError(ctx, errors.New("Project Not Null "))
		return
	}

	var isGit int32 = biz.GitProject
	if !p.IsGit {
		isGit = biz.NoneGitProject
	}
	info, err := DProject.Where(DProject.ID.Eq(p.ID)).Updates(model.Project{
		Name:       p.Name,
		Status:     p.Status,
		Branch:     &p.Branch,
		ProjectURL: &p.ProjectURL,
		PrivateKey: &p.PrivateKey,
		WorkDir:    &p.WorkDir,
		IsGit:      &isGit,
		Bin:        &p.Bin,
		Parameters: &p.Parameters,
	})
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if info.RowsAffected == 0 {
		handler.LaunchError(ctx, errors.New("Update Failed Can't Find ProjectInfo "))
		return
	}
	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Msg: "Update Ok",
	})
}

func GetProjectList(c context.Context, ctx *app.RequestContext) {
	projects, err := DProject.Order(DProject.CreatedAt.Desc()).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Data: projects,
	})
}

func GetProjectBuildStatus(c context.Context, ctx *app.RequestContext) {
	pid, err := strconv.Atoi(ctx.Param(biz.PID))
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	p, err := DProject.Select(DProject.Status).Where(DProject.ID.Eq(int32(pid))).Find()
	if err != nil {
		handler.LaunchError(ctx, err)
		return
	}
	if len(p) == 0 {
		handler.LaunchError(ctx, errors.New("Record Not Found "))
		return
	}

	ctx.JSON(consts.StatusOK, &biz.JsonModel{
		Data: p[0].Status,
	})
}
