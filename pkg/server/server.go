// Package server provides two server types: Server and
// HTTPTestServer.
//
// Server represents the server side implementation of Git. It wraps a
// git repository and the required sessions for responding to
// git-upload-pack and git-receive-pack requests after discovering
// available references. HTTPTestServer is an additional server layer
// which provides a convenient way to use the library as a test
// service over HTTP.
//
// You can use your own HTTP Server by passing your HTTP Request
// multiplexer to `SetupRoutes` if your mux implements `server.Router`
// interface which the `ServeMux` as well as `gorilla/mux` does.  If
// you are however using Gin, you can use `SetupGinRoutes` which
// accepts `gin.IRouter`.
package server

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	tsrv "github.com/go-git/go-git/v5/plumbing/transport/server"
)

const (
	defSessionTimeout = 5 * time.Minute
)

var (
	ErrRepoUninitialized = fmt.Errorf("git repo not initialized")
	ErrOwnerMissing      = fmt.Errorf("owner is empty")
	ErrRepoNameMissing   = fmt.Errorf("repoName is empty")
	ErrNilServer         = fmt.Errorf("server is nil")
	ErrInvalidAuth       = fmt.Errorf("invalid auth")
)

// BasicAuth is used to carry authentication for the HTTP endpoints.
type BasicAuth struct {
	Username string
	Password string
}

// Server holds the Git repository as well as the sessions for
// git-upload-pack and git-receive-pack operations.
type Server struct {
	Owner    string
	RepoName string

	SessionTimeout time.Duration
	basicAuth      BasicAuth

	repo *git.Repository

	// upSession represents git-upload-pack
	upSession transport.UploadPackSession
	// rpSession represents git-receive-pack
	rpSession transport.ReceivePackSession
}

type Option func(*Server)

func New(repo *git.Repository, owner, repoName string, opts ...Option) (*Server, error) {
	if repo == nil {
		return nil, ErrRepoUninitialized
	}

	if owner == "" {
		return nil, ErrOwnerMissing
	}

	if repoName == "" {
		return nil, ErrRepoNameMissing
	}

	_, err := repo.Reference(plumbing.HEAD, false)
	if err != nil {
		return nil, fmt.Errorf("git reference: %w", err)
	}

	srv := &Server{
		Owner:    strings.ToLower(owner),
		RepoName: strings.ToLower(repoName),

		SessionTimeout: defSessionTimeout,
		basicAuth: BasicAuth{
			Username: "",
			Password: "",
		},

		repo: repo,

		upSession: nil,
		rpSession: nil,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if err := srv.newSessions(); err != nil {
		return nil, fmt.Errorf("newSessions: %w", err)
	}

	return srv, nil
}

// Load provides the object store for the given end point to satisfy
// Go-Git sessions.
func (s *Server) Load(*transport.Endpoint) (storer.Storer, error) { //nolint:ireturn
	if s.repo == nil {
		return nil, ErrRepoUninitialized
	}

	return s.repo.Storer, nil
}

// RepoPath returns the relative path to the Git repository, it should
// be used together with base URL of the HTTP server.
func (s *Server) RepoPath() string {
	return path.Join(s.Owner, fmt.Sprintf("%s.git", s.RepoName))
}

func WithBasicAuth(ba BasicAuth) Option {
	return func(s *Server) {
		s.basicAuth = ba
	}
}

func (s *Server) newSessions() error {
	gitSrv := tsrv.NewServer(s)

	upSession, err := gitSrv.NewUploadPackSession(nil, nil)
	if err != nil {
		return fmt.Errorf("new UploadPackSession: %w", err)
	}

	rpSession, err := gitSrv.NewReceivePackSession(nil, nil)
	if err != nil {
		return fmt.Errorf("new ReceivePackSession: %w", err)
	}

	s.upSession = upSession
	s.rpSession = rpSession

	return nil
}

func (s *Server) authenticate(username, password string, _ bool) error {
	if s.basicAuth == (BasicAuth{Username: "", Password: ""}) {
		return nil
	}

	if s.basicAuth.Username == username &&
		s.basicAuth.Password == password {
		return nil
	}

	return ErrInvalidAuth
}
