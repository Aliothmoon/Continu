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

	// Project List
	g.GET("/projects", project.GetProjectList)

	//  Add Project
	g.POST("/project", project.AddProject)
	//  Update Project
	g.PUT("/project", project.UpdateProject)
	// Delete Project
	g.DELETE("/project/:PID", project.DelProject)

	// Get Project Status
	g.GET("/project/status/:PID", project.GetProjectBuildStatus)

	// Get Project Build List
	g.GET("/build/history/:PID", build.GetBuildHistoryList)

	// Build Project
	g.POST("/build/:PID", build.AddBuildTask)
	// Cancel Build Task
	g.POST("/build/cancel/:RID", build.CancelBuildTask)

	// Get Build Record Log
	g.GET("/log/:RID", logs.GetBuildLogs)

	// WebHooks
	g.POST("/hooks/:id", hooks.ProcessWebHooks)
}
