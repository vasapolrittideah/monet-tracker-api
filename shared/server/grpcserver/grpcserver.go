package grpcserver

import (
	"fmt"
	"net"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServiceRegistrar func(server *grpc.Server)

type GRPCServer struct {
	Addr             string
	ServiceRegistrar ServiceRegistrar
	server           *grpc.Server
}

func (s *GRPCServer) Start() error {
	s.server = grpc.NewServer()
	if s.ServiceRegistrar != nil {
		s.ServiceRegistrar(s.server)
	}

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.Addr, err)
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	log.Infof("🚀 gRPC server started on %s", s.Addr)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc server: %v", err)
	}

	return nil
}

func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()
	return nil
}
