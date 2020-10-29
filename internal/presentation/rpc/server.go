package rpc

import (
	"context"
	"log"

	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing/command"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/brycedarling/go-practical-microservices/internal/practicalpb"
)

// Server ...
type Server struct {
	practicalpb.PracticalServiceServer
	viewingQuery viewing.Query
	eventStore   eventstore.Store
}

var _ practicalpb.PracticalServiceServer = (*Server)(nil)

// NewServer ...
func NewServer(conf *config.Config) *Server {
	return &Server{nil, conf.ViewingQuery, conf.EventStore}
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
