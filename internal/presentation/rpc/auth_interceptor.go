package rpc

import (
	"context"

	"github.com/brycedarling/go-practical-microservices/internal/presentation/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor ...
type AuthInterceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

// NewAuthInterceptor ...
func NewAuthInterceptor() AuthInterceptor {
	return &authInterceptor{}
}

type authInterceptor struct{}

// Unary ...
func (ai *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, err := ai.authorize(ctx, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream ...
func (ai *authInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if _, err := ai.authorize(stream.Context(), info.FullMethod); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (ai *authInterceptor) authorize(ctx context.Context, method string) (*web.CustomClaims, error) {
	if method == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" ||
		method == "/practical.PracticalService/Login" {
		return nil, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authorization := md["authorization"]
	if len(authorization) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	claims, err := web.ParseJWT(authorization[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
	}

	return claims, nil
}
