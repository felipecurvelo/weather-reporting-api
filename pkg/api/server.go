package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Resource interface {
	Register(*httprouter.Router)
}

type Server struct {
	httpServer  *http.Server
	mainContext context.Context
	Router      *httprouter.Router
	stop        chan os.Signal
}

type ServerOptions struct {
	Port int
}

func (s *Server) GetHttpHandler() http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(s.mainContext)
		s.httpServer.Handler.ServeHTTP(w, r)
	})
	return handler
}

func (s *Server) Start() *Server {
	s.stop = make(chan os.Signal)
	signal.Notify(s.stop, syscall.SIGTERM, syscall.SIGINT)

	startServer := func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(fmt.Sprintf("[ListenAndServe]: %s\n", err))
			s.stop <- syscall.SIGQUIT
		}
	}

	go startServer()

	return s
}

func (s *Server) WaitForShutdownSignal() *Server {
	<-s.stop
	return s
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	return err
}

func (s *Server) RegisterResource(resource Resource) *Server {
	resource.Register(s.Router)
	return s
}

func NewServer(ctx context.Context, options *ServerOptions) *Server {
	server := &Server{
		httpServer:  &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", options.Port)},
		mainContext: ctx,
		Router:      httprouter.New(),
	}

	mainHandler := http.Handler(server.Router)
	contextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)
		mainHandler.ServeHTTP(w, r)
	})
	server.httpServer.Handler = contextHandler

	return server
}
