package rpc

import (
	"context"
	"log"
	"net"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	identityCommand "github.com/brycedarling/go-practical-microservices/internal/domain/identity/command"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	viewingCommand "github.com/brycedarling/go-practical-microservices/internal/domain/viewing/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/brycedarling/go-practical-microservices/internal/practicalpb"
	"github.com/brycedarling/go-practical-microservices/internal/presentation/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server ...
type Server struct {
	practicalpb.PracticalServiceServer
	grpcServer     *grpc.Server
	identityQuery  identity.Query
	viewingQuery   viewing.Query
	passwordHasher identity.PasswordHasher
	eventStore     eventstore.Store
	listener       net.Listener
	env            string
}

var _ practicalpb.PracticalServiceServer = (*Server)(nil)

// NewServer ...
func NewServer(conf *config.Config) (*Server, func(), error) {
	l, shutdownListener, err := web.NewTCPListener(":50051")
	if err != nil {
		return nil, nil, err
	}

	authInterceptor := NewAuthInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	s := &Server{
		env:            conf.Env.Env,
		eventStore:     conf.EventStore,
		identityQuery:  conf.IdentityQuery,
		viewingQuery:   conf.ViewingQuery,
		passwordHasher: conf.PasswordHasher,
		listener:       l,
		grpcServer:     grpcServer,
	}

	practicalpb.RegisterPracticalServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	return s, func() {
		shutdownListener()

		grpcServer.GracefulStop()
	}, nil
}

// Listen ...
func (s *Server) Listen() {
	log.Printf("Starting gRPC server in %s on %s", s.env, s.listener.Addr())
	if err := s.grpcServer.Serve(s.listener); err != nil && err != web.ErrShutdown {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Login ...
func (s *Server) Login(ctx context.Context, req *practicalpb.LoginRequest) (*practicalpb.LoginResponse, error) {
	log.Printf("Login invoked: %v", req)

	traceID := "TODO: ADD TRACEID TO GRPC CONTEXT"

	cmd, err := identityCommand.NewAuthenticateCommand(
		s.eventStore, s.identityQuery, s.passwordHasher, traceID, req.Email, req.Password,
	)
	if err != nil {
		return nil, err
	}

	id, err := cmd.Execute()
	if err != nil {
		if _, ok := err.(identityCommand.ErrAuthenticationFailed); ok {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}
		log.Println("Unexpected error authenticating:", err)
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	signedToken, err := web.SignJWT(id.UserID)
	if err != nil {
		return nil, err
	}

	return &practicalpb.LoginResponse{AccessToken: signedToken}, nil
}

// Viewing ...
func (s *Server) Viewing(ctx context.Context, req *practicalpb.ViewingRequest) (*practicalpb.ViewingResponse, error) {
	log.Printf("Viewing invoked: %v", req)
	viewing, err := s.viewingQuery.Find()
	if err != nil {
		return nil, err
	}
	return &practicalpb.ViewingResponse{
		Viewing: &practicalpb.Viewing{
			VideosWatched: int32(viewing.VideosWatched),
		},
	}, nil
}

// RecordViewing ...
func (s *Server) RecordViewing(ctx context.Context, req *practicalpb.RecordViewingRequest) (*practicalpb.RecordViewingResponse, error) {
	log.Printf("RecordViewing invoked: %v", req)

	traceID := ""
	userID := ""

	if traceID == "" {
	}

	cmd, err := viewingCommand.NewViewVideoCommand(s.eventStore, traceID, &userID, req.VideoId)
	if err != nil {
		return nil, err
	}

	err = cmd.Execute()
	if err != nil {
		return nil, err
	}

	return &practicalpb.RecordViewingResponse{}, nil
}
