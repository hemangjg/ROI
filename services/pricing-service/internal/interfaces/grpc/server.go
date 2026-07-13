package grpcapi

import (
	"fmt"
	"net"

	pricingv1 "github.com/ai-finops/ai-finops/packages/proto/gen/go/pricing/v1"
	"google.golang.org/grpc"
)

// Server wraps the gRPC listener and server.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewServer registers pricing gRPC services on the given port.
func NewServer(port int, pricingServer *PricingServer) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("listen grpc: %w", err)
	}

	srv := grpc.NewServer()
	pricingv1.RegisterPricingServiceServer(srv, pricingServer)

	return &Server{
		grpcServer: srv,
		listener:   lis,
	}, nil
}

// Serve blocks until the gRPC server stops.
func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.listener)
}

// Addr returns the bound listener address.
func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

// GracefulStop stops the gRPC server gracefully.
func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}