package grpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/seb7887/janus/internal/server/grpc/interceptor"
	"github.com/seb7887/janus/janusrpc"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	Serve(ctx context.Context) error
}

type grpcServer struct {
	grpcAddr string
}

func New(addr string) GRPCServer {
	return &grpcServer{
		grpcAddr: addr,
	}
}

func (s *grpcServer) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(withUnaryInterceptor())
	serviceServer := NewJanusGRPCServer()

	janusrpc.RegisterJanusServiceServer(grpcServer, serviceServer)

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func withUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptor.AuthorizationInterceptor,
	))
}
