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
	Server           *grpc.Server
	ServiceRegistrar ServiceRegistrar
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.Addr, err)
	}

	s.Server = grpc.NewServer()
	if s.ServiceRegistrar != nil {
		s.ServiceRegistrar(s.Server)
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.Server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	log.Infof("🚀 gRPC server started on %s", s.Addr)

	if err := s.Server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc server: %v", err)
	}

	return nil
}

func (s *GRPCServer) Stop() error {
	s.Server.GracefulStop()
	return nil
}
