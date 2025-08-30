package auth

import (
	"context"

	ssov1 "github.com/TimNikolaev/grpc-protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) SignUp(context.Context, *ssov1.SignUpRequest) (*ssov1.SignUpResponse, error) {
	return &ssov1.SignUpResponse{
		UserId: 1,
	}, nil
}

func (s *serverAPI) SignIn(context.Context, *ssov1.SignInRequest) (*ssov1.SignInResponse, error) {
	panic("implement me")
}
