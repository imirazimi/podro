package adapter

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HTTPServer struct {
	Config HTTPServerConfig
	App    *fiber.App
}

type HTTPServerConfig struct {
	Port uint `koanf:"port"`
}

func NewHTTPServer(config HTTPServerConfig) *HTTPServer {
	return &HTTPServer{App: fiber.New(), Config: config}
}

func (s *HTTPServer) Start() error {
	return s.App.Listen(fmt.Sprintf(":%d", s.Config.Port))
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.App.ShutdownWithContext(ctx)
}
