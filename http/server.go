package http

import (
	app "github.com/leartgjoni/go-rest-template"
	"net"
	"net/http"
)

type Server struct {
	ln net.Listener

	// Services
	UserService app.UserService
	ArticleService app.ArticleService

	// Handlers
	authHandler *AuthHandler
	articleHandler *ArticleHandler

	// Server options.
	Addr        string // bind address
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Open() error {
	s.initializeHandlers()

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	go http.Serve(s.ln, s.router())

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// initialize handlers server needs
func (s *Server) initializeHandlers() {
	s.authHandler = NewAuthHandler(s.UserService)
	s.articleHandler = NewArticleHandler(s.ArticleService)
}

// handlePing handles health check from kubernetes.
func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("healthy"))
}
