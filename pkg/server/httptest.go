package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-git/go-git/v5"
)

type HTTPTestServer struct {
	Server *Server
	TS     *httptest.Server
}

// NewHTTPTest initialises a new Git Server as well as a HTTP test
// server which is started.
func NewHTTPTest(repo *git.Repository, owner, repoName string) (*HTTPTestServer, error) {
	server, err := New(repo, owner, repoName)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	ts := httptest.NewServer(mux)

	h := &HTTPTestServer{
		Server: server,
		TS:     ts,
	}

	return h, nil
}

// URL returns the full path to the repository, it is the $GIT_URL
// reference in the Git HTTP protocol
// https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt.
func (h *HTTPTestServer) URL() string {
	return fmt.Sprintf("%s/%s", h.TS.URL, h.Server.RepoPath())
}

func (h *HTTPTestServer) Stop() {
	h.TS.Close()
}
