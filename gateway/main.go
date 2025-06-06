package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	proto "github.com/vasapolrittideah/money-tracker-api/protogen"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := bootstrap.App()
	defer app.Close()

	ropts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:     true,
				EmitDefaultValues: true,
			},
		}),
	}

	mux := runtime.NewServeMux(ropts...)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	)}

	// user
	userEndpoint := fmt.Sprintf("%v:%v", app.Config.Server.UserServiceHost, app.Config.Server.UserServicePort)
	err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userEndpoint, opts)
	if err != nil {
		log.Fatal("failed to register user service: %v", err)
	}

	log.Infof("🚀 server started on port %v", app.Config.Server.GatewayPort)

	addr := fmt.Sprintf(":%v", app.Config.Server.GatewayPort)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("failed to start server: %v", err)
	}
}
