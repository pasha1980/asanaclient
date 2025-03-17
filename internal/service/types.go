package service

import "context"

type AsanaClient interface {
	FetchUsers(context.Context, FetchInput) (FetchOutput[User], error)
	FetchProjects(context.Context, FetchInput) (FetchOutput[Project], error)
	FetchWorkspaces(context.Context) ([]Workspace, error)
}

type Storage interface {
	Save(context.Context, string, any) error
}

type FetchInput struct {
	Limit     int
	Workspace string
	Offset    *string
}

type FetchOutput[T any] struct {
	Data       []T
	NextOffset string
}

type User struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Project struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Workspace struct {
	GID  string
	Name string
}

type ExtractorMode uint

const (
	ExtractorModeFiveMinutes = iota
	ExtractorModeThirtySeconds
)
