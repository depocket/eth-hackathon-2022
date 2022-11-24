package service

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/repo"
	"go.uber.org/zap"
)

type AddressService struct {
	log  *zap.Logger
	repo repo.DgraphInterface
}

func NewAddressService(log *zap.Logger, repo repo.DgraphInterface) *AddressService {
	return &AddressService{
		log:  log,
		repo: repo,
	}
}

type AddressInterface interface {
	FullFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error)
	InFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error)
	OutFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error)
	Path(ctx context.Context, req model.PathRequest) (interface{}, error)
}

func (s *AddressService) FullFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error) {
	return s.repo.FullFlow(ctx, req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) InFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error) {
	return s.repo.InFlow(ctx, req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) OutFlow(ctx context.Context, req model.FlowRequest) (*model.ResponseFlow, error) {
	return s.repo.OutFlow(ctx, req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) Path(ctx context.Context, req model.PathRequest) (interface{}, error) {
	return s.repo.Path(ctx, req.Path, req.FromAddress, req.ToAddress)
}
