package server

import (
	"path"

	"github.com/gin-gonic/gin"
)

const (
	infoRefs    = "info/refs"
	uploadPack  = "git-upload-pack"
	receivePack = "git-receive-pack"
)

func (s *Server) pathInfoRefs() string {
	return path.Join("/", s.RepoPath(), infoRefs)
}

func (s *Server) pathUploadPack() string {
	return path.Join("/", s.RepoPath(), uploadPack)
}

func (s *Server) pathReceivePack() string {
	return path.Join("/", s.RepoPath(), receivePack)
}

// SetupRoutes adds required Git HTTP handlers to provided request
// multiplexer.
func (s *Server) SetupRoutes(r Router) {
	r.HandleFunc(s.pathInfoRefs(), s.GetInfoRefs)
	r.HandleFunc(s.pathUploadPack(), s.GetUploadPack)
	r.HandleFunc(s.pathReceivePack(), s.GetReceivePack)
}

// SetupGinRoutes adds required Git HTTP handlers to provided Gin
// IRouter, as Gin has a different interface we will wrap it to make
// the library easier to use for Gin users.
func (s *Server) SetupGinRoutes(ginRouter gin.IRouter) {
	ginRouter.GET(s.pathInfoRefs(), func(c *gin.Context) {
		s.GetInfoRefs(c.Writer, c.Request)
	})

	ginRouter.POST(s.pathUploadPack(), func(c *gin.Context) {
		s.GetUploadPack(c.Writer, c.Request)
	})

	ginRouter.POST(s.pathReceivePack(), func(c *gin.Context) {
		s.GetReceivePack(c.Writer, c.Request)
	})
}
