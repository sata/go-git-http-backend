package server_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
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
