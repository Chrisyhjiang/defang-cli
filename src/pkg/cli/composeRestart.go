package cli

import (
	"context"

	"github.com/DefangLabs/defang/src/pkg/cli/client"
	"github.com/DefangLabs/defang/src/pkg/types"
)

func ComposeRestart(ctx context.Context, client client.Client) (types.ETag, error) {
	project, err := client.LoadProject(ctx)
	if err != nil {
		return "", err
	}
	names := make([]string, 0, len(project.Services))
	for _, service := range project.Services {
		names = append(names, NormalizeServiceName(service.Name))
	}

	return Restart(ctx, client, names...)
}
