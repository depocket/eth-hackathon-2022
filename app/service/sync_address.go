package service

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/repo"
	"depocket.io/app/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/big"
)

type SyncAddressService struct {
	log  *zap.Logger
	repo repo.DgraphInterface
}

func NewSyncAddressService(log *zap.Logger, repo repo.DgraphInterface) *SyncAddressService {
	return &SyncAddressService{
		log:  log,
		repo: repo,
	}
}

type SyncAddressInterface interface {
	SyncAddress(ctx context.Context, txns []model.Transaction, symbols map[string]model.Token) error
}

func (s *SyncAddressService) SyncAddress(ctx context.Context, txns []model.Transaction, symbols map[string]model.Token) error {

	for _, txn := range txns {
		action := model.TransferAction{}
		if err := json.Unmarshal(txn.DecodedInput.RawMessage, &action); err != nil {
			s.log.Sugar().Error(err)
			continue
		}
		_, err := s.repo.GetByTransaction(ctx, txn.Hash)
		if err != nil && !utils.IsNotFound(err) {
			s.log.Sugar().Error(err)
			return err
		}
		if utils.IsNotFound(err) {
			// sender
			uidSender, err := s.repo.GetByAddress(ctx, txn.FromAddress)
			if err != nil && !utils.IsNotFound(err) {
				s.log.Sugar().Error(err)
				return err
			}

			uidRecipient, err := s.repo.GetByAddress(ctx, action.Recipient)
			if err != nil && !utils.IsNotFound(err) {
				s.log.Sugar().Error(err)
				return err
			}

			alias := fmt.Sprintf("%s -> %v", symbols[txn.ToAddress].Symbol, utils.ConvertBalance(&action.Amount, big.NewInt(int64(symbols[txn.ToAddress].Decimals))))
			if err := s.repo.CreateNode(ctx, "txn_id", action.Recipient, &model.TransactionDgraph{
				Amount: float64(action.Amount.Int64()),
				Recipient: model.AddressDgraph{
					UID:     uidRecipient,
					Address: action.Recipient,
					Name:    action.Recipient,
					Type:    "test-reci",
					DType:   []string{"Address"},
				},
				Sender: model.AddressDgraph{
					UID:     uidSender,
					Address: txn.FromAddress,
					Name:    txn.FromAddress,
					Type:    "test-send",
					DType:   []string{"Address"},
				},
				Name:         alias,
				TokenAddress: txn.ToAddress,
				TxnId:        txn.Hash,
				TxnTime:      txn.Block.Time,
				DType:        []string{"Transaction"},
			}); err != nil {
				s.log.Sugar().Error(err)
				return err
			}
		}
	}
	return nil
}
