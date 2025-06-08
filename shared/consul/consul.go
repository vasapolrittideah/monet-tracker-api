package consul

import (
	"fmt"
	"time"

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

func (c *ConsulClient) RegisterService(
	serviceID,
	serviceName,
	serviceAddress string,
	servicePort int,
	checkInterval time.Duration,
	deregisterAfter time.Duration,
) error {
	registration := &consulapi.AgentServiceRegistration{
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

	return c.client.Agent().ServiceRegister(registration)
}

func (c *ConsulClient) DeregisterService(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *ConsulClient) DiscoverServices(serviceName string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (c *ConsulClient) ConnectToService(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	services, err := c.DiscoverServices(serviceName)
	if err != nil {
		return nil, fmt.Errorf("no healthy instances of service %s found: %v", serviceName, err)
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
