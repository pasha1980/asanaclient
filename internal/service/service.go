package service

import (
	"context"
	"fmt"
	"github.com/pasha1980/asanaclient/config"
	"log/slog"
	"sync"
	"time"
)

func RunExtractor(ctx context.Context, mode ExtractorMode, client AsanaClient, storage Storage) {
	tickDuration := defineTickerDuration(mode)
	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case <-ticker.C:
			err := Extract(ctx, client, storage)
			if err != nil {
				handleError(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func Extract(ctx context.Context, client AsanaClient, storage Storage) error {
	slog.Info("------------------------")
	slog.Info("Extraction initiated")
	defer slog.Info("Extraction completed")

	workspaces, err := client.FetchWorkspaces(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for i := range workspaces {
		wg.Add(2)

		go func() {
			err = extractUsers(ctx, client, storage, workspaces[i].GID)
			if err != nil {
				handleError(err)
			}
			wg.Done()
		}()

		go func() {
			err = extractProjects(ctx, client, storage, workspaces[i].GID)
			if err != nil {
				handleError(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func extractUsers(ctx context.Context, client AsanaClient, storage Storage, workspaceID string) error {
	slog.Info("User extraction initiated")
	defer slog.Info("User extraction completed")

	var offset string
	cfg := config.Get()
	for {
		input := FetchInput{
			Workspace: workspaceID,
			Limit:     cfg.AsanaLimit,
		}

		if len(offset) > 0 {
			input.Offset = &offset
		}

		users, err := client.FetchUsers(ctx, input)
		if err != nil {
			return err
		}
		if len(users.Data) == 0 {
			return nil
		}

		offset = users.NextOffset

		for _, user := range users.Data {
			err = storage.Save(ctx, fmt.Sprintf("user.%s", user.GID), user)
			if err != nil {
				return err
			}
		}

		if len(users.NextOffset) == 0 {
			return nil
		}
	}
}

func extractProjects(ctx context.Context, client AsanaClient, storage Storage, workspaceID string) error {
	slog.Info("Projects extraction initiated")
	defer slog.Info("Projects extraction completed")

	var offset string
	cfg := config.Get()
	for {
		input := FetchInput{
			Workspace: workspaceID,
			Limit:     cfg.AsanaLimit,
		}

		if len(offset) > 0 {
			input.Offset = &offset
		}

		projects, err := client.FetchProjects(ctx, input)
		if err != nil {
			return err
		}
		if len(projects.Data) == 0 {
			return nil
		}

		offset = projects.NextOffset

		for _, user := range projects.Data {
			err = storage.Save(ctx, fmt.Sprintf("project.%s", user.GID), user)
			if err != nil {
				return err
			}
		}

		if len(projects.NextOffset) == 0 {
			return nil
		}
	}
}

func handleError(err error) {
	slog.Error(err.Error())
}

func defineTickerDuration(mode ExtractorMode) time.Duration {
	switch mode {
	case ExtractorModeFiveMinutes:
		return 5 * time.Minute
	default:
		return 30 * time.Second
	}
}
