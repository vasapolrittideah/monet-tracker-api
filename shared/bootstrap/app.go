package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/consul"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"github.com/vasapolrittideah/money-tracker-api/shared/server"
	"github.com/vasapolrittideah/money-tracker-api/shared/server/grpcserver"
	"github.com/vasapolrittideah/money-tracker-api/shared/server/httpserver"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	name           string
	config         *config.Config
	db             *gorm.DB
	servers        []server.Server
	clients        *consul.ClientRegistry
	deregisterFunc func() error
}

func NewApp(name string, config *config.Config, db *gorm.DB) *App {
	return &App{
		name:   name,
		config: config,
		db:     db,
	}
}

func (a *App) Run() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, srv := range a.servers {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				log.Errorf("failed to start %s server: %v", a.name, err)
			}
		}(srv)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Infof("received signal %v, shutting down...", sig)
	case <-ctx.Done():
		log.Info("context cancelled, shutting down...")
	}

	a.Stop()

	wg.Wait()
	log.Infof("👋 %s service stopped gracefully", a.name)
}

func (a *App) Stop() {
	database.Close(a.db)

	for _, s := range a.servers {
		_ = s.Stop()
	}

	if a.deregisterFunc != nil {
		_ = a.deregisterFunc()
	}
}

func (a *App) GetGRPCClient(name string) *grpc.ClientConn {
	if a.clients == nil {
		log.Fatal("client registry is not initialized")
	}

	return a.clients.Get(name)
}

func (a *App) AddHTTPServer(addr string, registrar httpserver.RouteRegistrar) {
	a.servers = append(a.servers, &httpserver.HTTPServer{
		Addr:           addr,
		RouteRegistrar: registrar,
	})
}

func (a *App) AddGRPCServer(addr string, registrar grpcserver.ServiceRegistrar) {
	a.servers = append(a.servers, &grpcserver.GRPCServer{
		Addr:             addr,
		ServiceRegistrar: registrar,
	})
}

func (a *App) ConnectGRPCClientsFromConsul(consulHost string, consulPort string, services []string) {
	conns, err := consul.ConnectGRPCServiceWithConsul(consulHost, consulPort, services)
	if err != nil {
		log.Fatalf("failed to connect to gRPC services via Consul: %v", err)
	}

	a.clients = consul.NewClientRegistry(conns)
}

func (a *App) RegisterGRPCServiceWithConsul(
	serviceID string,
	serviceName string,
	serviceHost string,
	servicePort string,
	checkInterval time.Duration,
	deregisterAfter time.Duration,
) {
	deregisterFunc, err := consul.RegisterGRPCServiceWithConsul(
		a.config.Server.ConsulHost,
		a.config.Server.ConsulPort,
		serviceID,
		serviceName,
		serviceHost,
		servicePort,
		checkInterval,
		deregisterAfter,
	)
	if err != nil {
		log.Fatalf("failed to register service to consul: %v", err)
	}
	a.deregisterFunc = deregisterFunc
}
