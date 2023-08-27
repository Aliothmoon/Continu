package router

import (
	"github.com/Aliothmoon/Continu/internal/web/handler/build"
	"github.com/Aliothmoon/Continu/internal/web/handler/hooks"
	"github.com/Aliothmoon/Continu/internal/web/handler/logs"
	"github.com/Aliothmoon/Continu/internal/web/handler/project"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(h *server.Hertz) {
	g := h.Group("/api")

	// ProjectInfo List
	g.GET("/projects", project.GetProjectList)

	//  Add ProjectInfo
	g.POST("/project", project.AddProject)
	//  Update ProjectInfo
	g.PUT("/project", project.UpdateProject)
	// Delete ProjectInfo
	g.DELETE("/project/:PID", project.DelProject)

	// Get ProjectInfo Status
	g.GET("/project/status/:PID", project.GetProjectBuildStatus)

	// Get ProjectInfo Build List
	g.GET("/build/history/:PID", build.GetBuildHistoryList)

	// Build ProjectInfo
	g.PUT("/build/:PID", build.AddBuildTask)

	// Cancel Build Task
	g.POST("/build/cancel/:RID", build.CancelBuildTask)

	// Get Build Record Log
	g.GET("/log/:RID", logs.GetBuildLogs)

	// WebHooks
	g.POST("/hooks/:PID", hooks.ProcessWebHooks)
}
