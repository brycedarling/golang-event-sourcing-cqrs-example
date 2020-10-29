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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server ...
type Server struct {
	practicalpb.PracticalServiceServer
	viewingQuery viewing.Query
	eventStore   eventstore.Store
	env          string
}

var _ practicalpb.PracticalServiceServer = (*Server)(nil)

// NewServer ...
func NewServer(conf *config.Config) *Server {
	return &Server{nil, conf.ViewingQuery, conf.EventStore, conf.Env.Env}
}

// Listen ...
func (s *Server) Listen() {
	grpcServer := grpc.NewServer()
	practicalpb.RegisterPracticalServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server in %s on %s", s.env, lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
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
