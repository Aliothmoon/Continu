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
	g.PUT("/project/:id", project.UpdateProject)
	// Delete Project
	g.DELETE("/project/:id", project.DelProject)

	// Build Project
	g.POST("/build", build.AddBuildTask)
	// Cancel Build Task
	g.POST("/build/cancel", build.CancelBuildTask)
	// Get Build List
	g.GET("/build/history/:id", build.GetBuildHistoryList)

	// Get Build Record Log
	g.GET("/log/:id", logs.GetBuildLogs)

	// WebHooks
	g.POST("/hooks/:id", hooks.ProcessWebHooks)
}
