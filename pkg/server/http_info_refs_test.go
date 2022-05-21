package server_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/sata-form3/go-git-http-backend/pkg/server"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // https://github.com/kunwardeep/paralleltest/issues/12
func TestInfoRefsWhenIncorrect(t *testing.T) {
	t.Parallel()

	service := "service"
	tests := map[string]struct {
		method         string
		queryParams    url.Values
		expectedStatus int
	}{
		"service name": {
			method:         http.MethodGet,
			queryParams:    url.Values{service: []string{"non-existent"}},
			expectedStatus: http.StatusForbidden,
		},
		"method": {
			method:         http.MethodPost,
			queryParams:    url.Values{service: []string{transport.ReceivePackServiceName}},
			expectedStatus: http.StatusBadRequest,
		},
		"too many params": {
			method:         http.MethodPost,
			queryParams:    url.Values{service: []string{"foo", transport.ReceivePackServiceName}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testRepo := emptyRepository(t)
			srv, err := server.NewHTTPTest(testRepo, owner, repoName)
			require.NoError(t, err, "server.New")

			client := &http.Client{
				Timeout: 1 * time.Second,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(
				ctx, test.method, fmt.Sprintf("%s/info/refs", srv.URL()), nil)
			require.NoError(t, err)

			req.URL.RawQuery = test.queryParams.Encode()

			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, test.expectedStatus, resp.StatusCode)
		})
	}
}
