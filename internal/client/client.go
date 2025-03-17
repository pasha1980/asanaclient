package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/pasha1980/asanaclient/config"
	"github.com/pasha1980/asanaclient/internal/service"
	"strconv"
)

type asanaClient struct {
	restyClient *resty.Client
}

func (c *asanaClient) FetchUsers(ctx context.Context, input service.FetchInput) (service.FetchOutput[service.User], error) {
	var result service.FetchOutput[service.User]

	req := c.restyClient.R().
		SetResult(&BaseAsanaResponse[service.User]{}).
		SetQueryParam("limit", strconv.Itoa(input.Limit)).
		SetQueryParam("workspace", input.Workspace)

	if input.Offset != nil {
		req.SetQueryParam("offset", *input.Offset)
	}

	resp, err := req.Get("/users")
	if err != nil {
		return result, err
	}

	parsedResp := resp.Result().(*BaseAsanaResponse[service.User])
	if parsedResp.NextPage != nil {
		result.NextOffset = parsedResp.NextPage.Offset
	}
	result.Data = parsedResp.Data
	return result, nil
}

func (c *asanaClient) FetchProjects(ctx context.Context, input service.FetchInput) (service.FetchOutput[service.Project], error) {
	var result service.FetchOutput[service.Project]

	req := c.restyClient.R().
		SetResult(&BaseAsanaResponse[service.Project]{}).
		SetQueryParam("limit", strconv.Itoa(input.Limit)).
		SetQueryParam("workspace", input.Workspace)

	if input.Offset != nil {
		req.SetQueryParam("offset", *input.Offset)
	}

	resp, err := req.Get("/projects")
	if err != nil {
		return result, err
	}

	parsedResp := resp.Result().(*BaseAsanaResponse[service.Project])
	if parsedResp.NextPage != nil {
		result.NextOffset = parsedResp.NextPage.Offset
	}
	result.Data = parsedResp.Data
	return result, nil
}

func (c *asanaClient) FetchWorkspaces(ctx context.Context) ([]service.Workspace, error) {
	req := c.restyClient.R().
		SetResult(&BaseAsanaResponse[service.Workspace]{})

	resp, err := req.Get("/workspaces")
	if err != nil {
		return nil, err
	}

	parsedResp := resp.Result().(*BaseAsanaResponse[service.Workspace])
	return parsedResp.Data, nil
}

func NewAsanaClient() service.AsanaClient {
	cfg := config.Get()

	restyClient := resty.New()
	restyClient.SetBaseURL(cfg.AsanaBaseUrl)
	restyClient.SetAuthToken(cfg.AsanaAccessToken)

	return &asanaClient{
		restyClient: restyClient,
	}
}
