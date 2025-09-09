package auth

import (
	"context"

	ssov1 "github.com/TimNikolaev/grpc-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	SignUpNewUser(ctx context.Context, email string, pass string) (uint64, error)
	SignIn(ctx context.Context, email string, pass string) (string, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyValue = 0
)

func (s *serverAPI) SignUp(ctx context.Context, req *ssov1.SignUpRequest) (*ssov1.SignUpResponse, error) {
	if err := s.validateSignUp(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.SignUpNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//TODO:
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignUpResponse{UserId: userID}, nil
}

func (s *serverAPI) SignIn(ctx context.Context, req *ssov1.SignInRequest) (*ssov1.SignInResponse, error) {
	if err := s.validateSignIn(req); err != nil {
		return nil, err
	}

	token, err := s.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//TODO
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignInResponse{Token: token}, nil
}

func (s *serverAPI) validateSignUp(req *ssov1.SignUpRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func (s *serverAPI) validateSignIn(req *ssov1.SignInRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}
