package cli

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/defang-io/defang/src/pkg/cli/client"
	defangv1 "github.com/defang-io/defang/src/protos/io/defang/v1"
)

func TestComposeStart(t *testing.T) {
	DoDryRun = true
	defer func() { DoDryRun = false }()

	project, err := LoadComposeWithProjectName("../../tests/testproj/compose.yaml", "tenant-id")
	if err != nil {
		t.Fatalf("LoadComposeWithProjectName() failed: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	_, err = ComposeStart(context.Background(), client.MockClient{UploadUrl: server.URL + "/"}, project, false)
	if !errors.Is(err, ErrDryRun) {
		t.Fatalf("ComposeStart() failed: %v", err)
	}
}

func TestComposeFixupEnv(t *testing.T) {
	project, err := LoadComposeWithProjectName("../../tests/fixupenv/compose.yaml", "tenant-id")
	if err != nil {
		t.Fatalf("LoadComposeWithProjectName() failed: %v", err)
	}

	services, err := convertServices(context.Background(), client.MockClient{}, project.Services, false)
	if err != nil {
		t.Fatalf("convertServices() failed: %v", err)
	}
	ui := slices.IndexFunc(services, func(s *defangv1.Service) bool { return s.Name == "ui" })
	if ui < 0 {
		t.Fatalf("convertServices() failed: expected service named 'ui'")
	}

	const expected = "http://mistral:8000"
	got := services[ui].Environment["API_URL"]
	if got != expected {
		t.Errorf("convertServices() failed: expected API_URL=%s, got %s", expected, got)
	}
}
