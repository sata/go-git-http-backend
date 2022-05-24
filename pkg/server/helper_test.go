package server_test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitfs "github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sata-form3/go-git-http-backend/pkg/server"
	"github.com/stretchr/testify/require"
)

func readFile(t *testing.T, fs billy.Filesystem, path string) string {
	t.Helper()

	file, err := fs.Open(path)
	require.NoError(t, err, "open file")

	bytes, err := ioutil.ReadAll(file)
	require.NoError(t, err, "read all")

	return string(bytes)
}

func emptyRepository(t *testing.T) *git.Repository {
	t.Helper()

	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	return repo
}

func repoWithInitCommit(t *testing.T, name, content string) *git.Repository {
	t.Helper()

	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	fs := worktree.Filesystem
	file, err := fs.Create(name)
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	err = worktree.AddGlob("*")
	require.NoError(t, err)

	_, err = worktree.Commit("initial commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "bob the builder",
			Email: "bob@builder.test",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	return repo
}

type cloneAssert struct {
	t *testing.T

	url string

	auth server.BasicAuth
}

type cloneAssertOption func(*cloneAssert)

func newCloneAssert(t *testing.T, url string, opts ...cloneAssertOption) *cloneAssert {
	t.Helper()

	cloneAssert := &cloneAssert{
		t: t,

		url: url,

		auth: server.BasicAuth{
			Username: "",
			Password: "",
		},
	}

	for _, opt := range opts {
		opt(cloneAssert)
	}

	return cloneAssert
}

func withAuth(auth server.BasicAuth) cloneAssertOption {
	return func(c *cloneAssert) {
		c.auth = auth
	}
}

func (c *cloneAssert) assert(filename, content string) {
	c.t.Helper()

	storage := gitfs.NewStorage(memfs.New(), cache.NewObjectLRUDefault())
	opts := &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: c.auth.Username,
			Password: c.auth.Password,
		},
		URL:           c.url,
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName("master"),
		SingleBranch:  false,
		Depth:         0,
	}

	repo, err := git.CloneContext(context.Background(), storage, memfs.New(), opts)
	require.NoError(c.t, err)
	require.NotEmpty(c.t, repo)

	wt, err := repo.Worktree()
	require.NoError(c.t, err)

	fs := wt.Filesystem
	readContent := readFile(c.t, fs, filename)

	require.Equal(c.t, content, readContent, "content mismatch")
}
