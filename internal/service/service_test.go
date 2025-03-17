package service_test

import (
	"context"
	"github.com/pasha1980/asanaclient/internal/service"
	"github.com/pasha1980/asanaclient/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestExtract_whenEmptyWorkspaces(t *testing.T) {
	asanaClient := mocks.NewServiceAsanaClient(t)
	asanaClient.EXPECT().FetchWorkspaces(mock.Anything).Return([]service.Workspace{}, nil)

	storage := mocks.NewServiceStorage(t)
	err := service.Extract(context.Background(), asanaClient, storage)
	assert.NoError(t, err)

	asanaClient.AssertNotCalled(t, "FetchUsers", mock.Anything, mock.Anything)
	asanaClient.AssertNotCalled(t, "FetchProjects", mock.Anything, mock.Anything)
	storage.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything)
}

func TestExtract_whenEmptyUsersAndProjects(t *testing.T) {
	asanaClient := mocks.NewServiceAsanaClient(t)
	asanaClient.EXPECT().FetchWorkspaces(mock.Anything).Return([]service.Workspace{
		{
			GID:  "1",
			Name: "1",
		},
	}, nil)
	asanaClient.EXPECT().FetchUsers(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.User]{
		Data:       []service.User{},
		NextOffset: "",
	}, nil)

	asanaClient.EXPECT().FetchProjects(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.Project]{
		Data:       []service.Project{},
		NextOffset: "",
	}, nil)

	storage := mocks.NewServiceStorage(t)
	err := service.Extract(context.Background(), asanaClient, storage)
	assert.NoError(t, err)

	storage.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything)
}

func TestExtract_whenOnePage(t *testing.T) {
	asanaClient := mocks.NewServiceAsanaClient(t)
	asanaClient.EXPECT().FetchWorkspaces(mock.Anything).Return([]service.Workspace{
		{"1", "1"},
	}, nil)
	asanaClient.EXPECT().FetchUsers(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.User]{
		Data: []service.User{
			{"1", "1"},
			{"2", "2"},
		},
		NextOffset: "",
	}, nil)

	asanaClient.EXPECT().FetchProjects(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.Project]{
		Data: []service.Project{
			{"1", "1"},
			{"2", "2"},
		},
		NextOffset: "",
	}, nil)

	storage := mocks.NewServiceStorage(t)
	storage.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := service.Extract(context.Background(), asanaClient, storage)
	assert.NoError(t, err)

	storage.AssertCalled(t, "Save", mock.Anything, "user.1", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "user.2", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.1", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.2", mock.Anything)
}

func TestExtract_whenTwoPages(t *testing.T) {
	offset := "secondPage"

	asanaClient := mocks.NewServiceAsanaClient(t)
	asanaClient.EXPECT().FetchWorkspaces(mock.Anything).Return([]service.Workspace{
		{"1", "1"},
	}, nil)
	asanaClient.EXPECT().FetchUsers(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.User]{
		Data: []service.User{
			{"1", "1"},
			{"2", "2"},
		},
		NextOffset: offset,
	}, nil)
	asanaClient.EXPECT().FetchUsers(mock.Anything, service.FetchInput{
		Workspace: "1",
		Offset:    &offset,
	}).Return(service.FetchOutput[service.User]{
		Data: []service.User{
			{"3", "3"},
			{"4", "4"},
		},
		NextOffset: "",
	}, nil)

	asanaClient.EXPECT().FetchProjects(mock.Anything, service.FetchInput{
		Workspace: "1",
	}).Return(service.FetchOutput[service.Project]{
		Data: []service.Project{
			{"1", "1"},
			{"2", "2"},
		},
		NextOffset: offset,
	}, nil)
	asanaClient.EXPECT().FetchProjects(mock.Anything, service.FetchInput{
		Workspace: "1",
		Offset:    &offset,
	}).Return(service.FetchOutput[service.Project]{
		Data: []service.Project{
			{"3", "3"},
			{"4", "4"},
		},
		NextOffset: "",
	}, nil)

	storage := mocks.NewServiceStorage(t)
	storage.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := service.Extract(context.Background(), asanaClient, storage)
	assert.NoError(t, err)

	storage.AssertCalled(t, "Save", mock.Anything, "user.1", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "user.2", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "user.3", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "user.4", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.1", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.2", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.3", mock.Anything)
	storage.AssertCalled(t, "Save", mock.Anything, "project.4", mock.Anything)
}
