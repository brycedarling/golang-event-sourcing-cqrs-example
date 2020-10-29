package rpc

import (
	"context"
	"log"
	"net"

	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/brycedarling/go-practical-microservices/internal/practicalpb"
	"github.com/brycedarling/go-practical-microservices/internal/presentation/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server ...
type Server struct {
	practicalpb.PracticalServiceServer
	grpcServer   *grpc.Server
	viewingQuery viewing.Query
	eventStore   eventstore.Store
	listener     net.Listener
	env          string
}

var _ practicalpb.PracticalServiceServer = (*Server)(nil)

// NewServer ...
func NewServer(conf *config.Config) (*Server, func(), error) {
	l, shutdownListener, err := web.NewTCPListener(":50051")
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()

	s := &Server{
		env:          conf.Env.Env,
		eventStore:   conf.EventStore,
		viewingQuery: conf.ViewingQuery,
		listener:     l,
		grpcServer:   grpcServer,
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

	cmd, err := command.NewViewVideoCommand(s.eventStore, traceID, &userID, req.VideoId)
	if err != nil {
		return nil, err
	}

	err = cmd.Execute()
	if err != nil {
		return nil, err
	}

	return &practicalpb.RecordViewingResponse{}, nil
}
