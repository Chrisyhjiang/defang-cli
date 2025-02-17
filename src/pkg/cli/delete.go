package cli

import (
	"context"

	"github.com/defang-io/defang/src/pkg/cli/client"
	"github.com/defang-io/defang/src/pkg/term"
	defangv1 "github.com/defang-io/defang/src/protos/io/defang/v1"
)

// Deprecated: Use ComposeStop instead.
func Delete(ctx context.Context, client client.Client, names ...string) (client.ETag, error) {
	term.Debug(" - Deleting service", names)

	if DoDryRun {
		return "", ErrDryRun
	}

	resp, err := client.Delete(ctx, &defangv1.DeleteRequest{Names: names})
	if err != nil {
		return "", err
	}
	return resp.Etag, nil
}
