package ghttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DataWorkbench/glog"
)

// HTTPServer is the configuration of http server
type ServerConfig struct {
	Address      string        `json:"address"       yaml:"address"       env:"ADDRESS"                   validate:"required"`
	ReadTimeout  time.Duration `json:"read_timeout"  yaml:"read_timeout"  env:"READ_TIMEOUT,default=30s"  validate:"required"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT,default=30s" validate:"required"`
	IdleTimeout  time.Duration `json:"idle_timeout"  yaml:"idle_timeout"  env:"DLE_TIMEOUT,default=30s"   validate:"required"`
	ExitTimeout  time.Duration `json:"exit_timeout"  yaml:"exit_timeout"  env:"EXIT_TIMEOUT,default=5m"   validate:"required"`
}

type Server struct {
	lp  *glog.Logger
	cfg *ServerConfig
	std *http.Server
}

func NewServer(ctx context.Context, cfg *ServerConfig, handler http.Handler) *Server {
	return &Server{
		lp:  glog.FromContext(ctx),
		cfg: cfg,
		std: &http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			ErrorLog:     log.New(os.Stderr, "httpServer: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lmsgprefix),
		},
	}
}

func (s *Server) ListenAndServe() (err error) {
	s.lp.Info().String("http server listening", s.cfg.Address).Fire()

	err = s.std.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.lp.Error().Error("listen and serve http server error", err).Fire()
		return
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	s.lp.Info().Msg("http server exit listen").Fire()
	return
}

func (s *Server) Shutdown(ctx context.Context) (err error) {
	if s == nil {
		return
	}

	if s.cfg.ExitTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.cfg.ExitTimeout)
		defer cancel()
	}

	s.lp.Info().Msg("waiting for http server shutdown").Fire()
	if err = s.std.Shutdown(ctx); err != nil {
		s.lp.Error().Error("http server shutdown error", err).Fire()
		return
	}
	s.lp.Info().Msg("http server has been shutdown").Fire()
	return
}
