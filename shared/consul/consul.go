package consul

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectGRPCServiceWithConsul connects to the given list of gRPC services registered in Consul
func ConnectGRPCServiceWithConsul(
	consulHost string,
	consulPort string,
	serviceNames []string,
) (map[string]*grpc.ClientConn, error) {
	address := fmt.Sprintf("%v:%v", consulHost, consulPort)
	consulClient, err := newConsulClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	log.Infof("connecting to consul at %s", address)

	return consulClient.createClientConnections(serviceNames)
}

// RegisterGRPCServiceWithConsul registers the gRPC service in Consul and deregisters on context cancellation
func RegisterGRPCServiceWithConsul(
	consulHost string,
	consulPort string,
	serviceID string,
	serviceName string,
	serviceAddress string,
	servicePort string,
	checkInterval time.Duration,
	deregisterAfter time.Duration,
) (func() error, error) {
	address := fmt.Sprintf("%v:%v", consulHost, consulPort)
	consulClient, err := newConsulClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	if err := consulClient.registerService(
		serviceID,
		serviceName,
		serviceAddress,
		servicePort,
		checkInterval,
		deregisterAfter,
	); err != nil {
		return nil, fmt.Errorf("failed to register service: %w", err)
	}

	log.Infof("🎉 %s registered successfully with consul", serviceName)

	deregisterFunc := func() error {
		if err := consulClient.deregisterService(serviceID); err != nil {
			return fmt.Errorf("failed to deregister %s: %v", serviceName, err)
		}

		log.Infof("%s deregistered successfully", serviceName)

		return nil
	}

	return deregisterFunc, nil
}

type ConsulClient struct {
	client *consulapi.Client
}

func newConsulClient(address string) (*ConsulClient, error) {
	config := consulapi.DefaultConfig()
	config.Address = address

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &ConsulClient{client: client}, nil
}

func (c *ConsulClient) createClientConnections(serviceNames []string) (map[string]*grpc.ClientConn, error) {
	conns := make(map[string]*grpc.ClientConn)

	for _, name := range serviceNames {
		conn, err := c.connectToService(name)
		if err != nil {
			for svc, c := range conns {
				if cerr := c.Close(); cerr != nil {
					log.Errorf("failed to close connection to %v: %v", svc, cerr)
				}
			}

			return nil, fmt.Errorf("failed to connect to service %v: %v", name, err)
		}

		conns[name] = conn
	}

	return conns, nil
}

func (c *ConsulClient) registerService(
	serviceID,
	serviceName,
	serviceAddress string,
	servicePort string,
	checkInterval time.Duration,
	deregisterAfter time.Duration,
) error {
	servicePortInt, _ := strconv.Atoi(servicePort)
	reg := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePortInt,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceAddress, servicePortInt),
			Interval:                       checkInterval.String(),
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: deregisterAfter.String(),
		},
	}

	return c.client.Agent().ServiceRegister(reg)
}

func (c *ConsulClient) deregisterService(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *ConsulClient) connectToService(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover services for %s: %w", serviceName, err)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	svc := services[0]
	target := fmt.Sprintf("%s:%d", svc.Service.Address, svc.Service.Port)

	if len(opts) == 0 {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service %s: %v", serviceName, err)
	}

	return conn, nil
}
