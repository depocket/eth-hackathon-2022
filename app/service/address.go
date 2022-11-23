package service

import (
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
	FullFlow(req model.FlowRequest) (interface{}, error)
	InFlow(req model.FlowRequest) (interface{}, error)
	OutFlow(req model.FlowRequest) (interface{}, error)
	Path(req model.PathRequest) (interface{}, error)
}

func (s *AddressService) FullFlow(req model.FlowRequest) (interface{}, error) {
	return s.repo.FullFlow(req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) InFlow(req model.FlowRequest) (interface{}, error) {
	return s.repo.InFlow(req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) OutFlow(req model.FlowRequest) (interface{}, error) {
	return s.repo.OutFlow(req.Depth, req.Address, req.Token, req.From, req.To)
}

func (s *AddressService) Path(req model.PathRequest) (interface{}, error) {
	return s.repo.Path(req.Path, req.FromAddress, req.ToAddress)
}
