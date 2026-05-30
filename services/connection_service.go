package services

import (
	"context"
	"whatsapp-backend/dto"
	"whatsapp-backend/repositories"
)

type ConnectionService interface {
	Connect(ctx context.Context, req dto.ConnectRequest) error
}
type connectionService struct {
	connectionRepo repositories.ConnectionRepository
}

func NewConnectionService(connectionRepo repositories.ConnectionRepository) ConnectionService {
	return &connectionService{connectionRepo: connectionRepo}
}

func (s *connectionService) Connect(ctx context.Context, req dto.ConnectRequest) error {
	return s.connectionRepo.Connect(ctx)
}
