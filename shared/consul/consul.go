package consul

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConsulClient struct {
	client *consulapi.Client
}

func NewConsulClient(address string) (*ConsulClient, error) {
	config := consulapi.DefaultConfig()
	config.Address = address

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &ConsulClient{client: client}, nil
}

func (c *ConsulClient) CreateGRPCClients(serviceNames []string) (map[string]*grpc.ClientConn, error) {
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

func (c *ConsulClient) RegisterService(
	serviceID,
	serviceName,
	serviceAddress string,
	servicePort int,
	checkInterval time.Duration,
	deregisterAfter time.Duration,
) error {
	reg := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePort,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceAddress, servicePort),
			Interval:                       checkInterval.String(),
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: deregisterAfter.String(),
		},
	}

	if err := c.client.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("failed to register service %v: %v", serviceName, err)
	}

	return nil
}

func (c *ConsulClient) DeregisterService(serviceID string) error {
	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service %v: %v", serviceID, err)
	}

	return nil
}

func (c *ConsulClient) discoverServices(serviceName string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover services for %s: %w", serviceName, err)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	return services, nil
}

func (c *ConsulClient) connectToService(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	services, err := c.discoverServices(serviceName)
	if err != nil {
		return nil, err
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
