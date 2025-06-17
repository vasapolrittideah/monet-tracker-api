package consul

import "google.golang.org/grpc"

type ClientRegistry struct {
	clients map[string]*grpc.ClientConn
}

func NewClientRegistry(clients map[string]*grpc.ClientConn) *ClientRegistry {
	return &ClientRegistry{clients: clients}
}

func (r *ClientRegistry) Get(name string) *grpc.ClientConn {
	return r.clients[name]
}
