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
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitfs "github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/sata-form3/go-git-http-backend/pkg/server"
	"github.com/stretchr/testify/require"
)

const (
	owner    = "bob"
	repoName = "shed"
)

func TestCloneHTTP(t *testing.T) {
	t.Parallel()

	fileName := "some_file"
	content := "some_content"
	testRepo := repoWithInitCommit(t, fileName, content)

	srv, err := server.NewHTTPTest(testRepo, owner, repoName)
	require.NoError(t, err, "server.New")

	storage := gitfs.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	opts := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "some-test",
			Password: "some-pass",
		},
		URL:           srv.URL(),
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		SingleBranch:  false,
		Depth:         0,
	}

	repo, err := git.CloneContext(context.Background(), storage, memfs.New(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	fs := wt.Filesystem
	readContent := readFile(t, fs, fileName)

	require.Equal(t, content, readContent, "content mismatch")
}

func TestCloneExternalHTTPWithGin(t *testing.T) {
	t.Parallel()

	fileName := "some_file"
	content := "some_content"
	testRepo := repoWithInitCommit(t, fileName, content)

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

	storage := gitfs.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	opts := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "some-test",
			Password: "some-pass",
		},
		URL:           fmt.Sprintf("%s/%s/%s.git", httpsrv.URL, owner, repoName),
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		SingleBranch:  false,
		Depth:         0,
	}

	repo, err := git.CloneContext(context.Background(), storage, memfs.New(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	fs := wt.Filesystem
	readContent := readFile(t, fs, fileName)

	require.Equal(t, content, readContent, "content mismatch")
}
