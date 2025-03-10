package auth

import (
	"context"

	protos "github.com/SemenShakhray/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	protos.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	protos.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *protos.LoginRequest) (*protos.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *protos.IsAdminRequest) (*protos.IsAdminResponse, error) {
	panic("implement me")
}
