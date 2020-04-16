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
	httpServer *http.Server
	stop       chan os.Signal
	Router     *httprouter.Router
}

type ServerOptions struct {
	Port int
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

func NewServer(options *ServerOptions) *Server {
	server := &Server{
		httpServer: &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", options.Port)},
		Router:     httprouter.New(),
	}

	handler := http.Handler(server.Router)
	server.httpServer.Handler = handler

	return server
}
