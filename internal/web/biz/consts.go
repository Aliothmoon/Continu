package biz

// Project
const (
	ProjectIdle = 2<<iota - 2
	ProjectPending
)

// Build Record
const (
	BuildPending = 2<<iota - 2
	BuildSuccess
	BuildFailed
)

// Judge is Git Project
const (
	NoneGitProject = iota
	GitProject
)
