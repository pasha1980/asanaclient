package client

type BaseAsanaResponse[T any] struct {
	Data     []T `json:"data"`
	NextPage *AsanaNextPageResponse
}

type AsanaNextPageResponse struct {
	Offset string
}
