package server_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitfs "github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/sata-form3/go-git-http-backend/pkg/server"
	"github.com/stretchr/testify/require"
)

const (
	owner    = "bob"
	repoName = "shed"
	filename = "somefile"
	content  = "some content"
)

func TestCloneHTTP(t *testing.T) {
	t.Parallel()

	testRepo := repoWithInitCommit(t, filename, content)

	srv, err := server.NewHTTPTest(testRepo, owner, repoName)
	require.NoError(t, err, "server.New")

	a := newCloneAssert(t, srv.URL())
	a.assert(filename, content)
}

func TestCloneHTTPWithAuth(t *testing.T) {
	t.Parallel()

	testRepo := repoWithInitCommit(t, filename, content)

	auth := server.BasicAuth{
		Username: "godmode",
		Password: "IDKFA",
	}

	srv, err := server.NewHTTPTest(testRepo, owner, repoName, server.WithBasicAuth(auth))
	require.NoError(t, err, "server.New")

	a := newCloneAssert(t, srv.URL(), withAuth(auth))
	a.assert(filename, content)
}

func TestCloneHTTPWitInvalidAuth(t *testing.T) {
	t.Parallel()

	testRepo := repoWithInitCommit(t, filename, content)

	auth := server.BasicAuth{
		Username: "godmode",
		Password: "IDKFA",
	}

	srv, err := server.NewHTTPTest(testRepo, owner, repoName, server.WithBasicAuth(auth))
	require.NoError(t, err, "server.New")

	storage := gitfs.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	opts := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "incorrect",
			Password: "incorrect",
		},
		URL:           srv.URL(),
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		SingleBranch:  false,
		Depth:         0,
	}
	_, err = git.CloneContext(context.Background(), storage, memfs.New(), opts)
	require.ErrorIs(t, err, transport.ErrAuthenticationRequired)
}

func TestCloneExternalHTTPWithGin(t *testing.T) {
	t.Parallel()

	testRepo := repoWithInitCommit(t, filename, content)

	defVal := gin.Mode()

	t.Cleanup(func() {
		gin.SetMode(defVal)
	})

	gin.SetMode(gin.TestMode)

	srv, err := server.New(testRepo, owner, repoName)
	require.NoError(t, err, "server.New")

	g := gin.New()
	srv.SetupGinRoutes(g)
	httpsrv := httptest.NewServer(g)

	url := fmt.Sprintf("%s/%s", httpsrv.URL, srv.RepoPath())

	a := newCloneAssert(t, url)
	a.assert(filename, content)
}
