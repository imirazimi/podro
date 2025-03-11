package podroapp

import (
	"interview/adapter"
)

func (h *Handler) SetupHttp(server *adapter.HTTPServer) *Handler {
	podroGP := server.App.Group("/api/v1/podroapp")
	podroGP.Get("/orders/weakly-report", h.GetProivedersWeaklyReport)
	return h
}

func SetupHandler(svc *Service) *Handler {
	return NewHandler(svc)
}

func SetupService(SQLDB *adapter.SQLDB, HTTPServer *adapter.HTTPServer) *Service {
	svc := NewService(
		NewDB(SQLDB),
	)
	SetupHandler(svc).SetupHttp(HTTPServer)
	return svc
}

func SetupClient(svc *Service, otp *adapter.OTP) *Client {
	return NewClient(svc, otp)
}

func (s *Service) SetClient(Client *Client) *Service {
	s.providerClient = Client
	s.smsClient = Client
	return s
}
