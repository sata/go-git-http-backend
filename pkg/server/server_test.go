package server_test

import (
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/sata-form3/go-git-http-backend/pkg/server"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // https://github.com/kunwardeep/paralleltest/issues/12
func TestNewWithInvalidArgs(t *testing.T) {
	t.Parallel()

	repo := repoWithInitCommit(t, "file", "content")

	tests := map[string]struct {
		repo     *git.Repository
		owner    string
		repoName string
		exp      *server.Server
		err      error
	}{
		"NilRepo": {
			nil,
			"owner",
			"name",
			nil,
			server.ErrRepoUninitialized,
		},
		"MissingOwner": {
			repo,
			"",
			"name",
			nil,
			server.ErrOwnerMissing,
		},
		"MissingRepoName": {
			repo,
			"owner",
			"",
			nil,
			server.ErrRepoNameMissing,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := server.New(test.repo, test.owner, test.repoName)
			require.Nil(t, actual)
			require.ErrorIs(t, err, test.err)
		})
	}
}

// Load should never cause panic, even if it's invoked on an nil struct.
// This is to catch regression.
func TestNoPanicsOnLoadWhenNilStruct(t *testing.T) {
	t.Parallel()

	s := &server.Server{}

	require.NotPanics(t, func() {
		actual, err := s.Load(&transport.Endpoint{})

		require.Equal(t, server.ErrRepoUninitialized, err)
		require.Nil(t, actual)
	})
}

//nolint:paralleltest // https://github.com/kunwardeep/paralleltest/issues/12
func TestURLsAreNormalized(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		owner    string
		repoName string
		exp      string
	}{
		"OwnerIsCapitalised": {
			"OWNER",
			"somerepo",
			"owner/somerepo.git",
		},
		"RepoIsCapitalised": {
			"owner",
			"SOMEREPO",
			"owner/somerepo.git",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s, err := server.NewHTTPTest(repoWithInitCommit(t, "foo", "bar"), test.owner, test.repoName)
			require.NoError(t, err)

			require.True(t, strings.HasSuffix(s.URL(), test.exp))
		})
	}
}
