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
	SyncAddress(ctx context.Context, chain string, txns []model.Transaction, symbols map[string]model.Token) error
}

func (s *SyncAddressService) SyncAddress(ctx context.Context, chain string, txns []model.Transaction, symbols map[string]model.Token) error {
	for _, txn := range txns {
		ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
		action := model.TransferAction{}
		if err := json.Unmarshal(txn.DecodedInput.RawMessage, &action); err != nil {
			s.log.Sugar().Error(err)
			cancel()
			continue
		}
		_, err := s.repo.GetByTransaction(ctx, txn.Hash)
		if err != nil && !utils.IsNotFound(err) {
			s.log.Sugar().Error(err)
			cancel()
			continue
		}
		if utils.IsNotFound(err) {
			// sender
			uidSender, err := s.repo.GetByAddress(ctx, txn.FromAddress)
			if err != nil && !utils.IsNotFound(err) {
				s.log.Sugar().Error(err)
				cancel()
				return err
			}

			uidRecipient, err := s.repo.GetByAddress(ctx, action.Recipient)
			if err != nil && !utils.IsNotFound(err) {
				s.log.Sugar().Error(err)
				cancel()
				return err
			}

			amount := utils.ConvertBalance(&action.Amount, big.NewInt(int64(symbols[txn.ToAddress].Decimals)))
			alias := fmt.Sprintf("%s:%v", symbols[txn.ToAddress].Symbol, utils.ConvertBalance(&action.Amount, big.NewInt(int64(symbols[txn.ToAddress].Decimals))))
			if err := s.repo.CreateNode(ctx, "txn_id", action.Recipient, &model.TransactionDgraph{
				Amount: *amount,
				Recipient: model.AddressDgraph{
					UID:     uidRecipient,
					Address: action.Recipient,
					Name:    action.Recipient,
					Chain:   chain,
					DType:   []string{"Address"},
				},
				Sender: model.AddressDgraph{
					UID:     uidSender,
					Address: txn.FromAddress,
					Name:    txn.FromAddress,
					Chain:   chain,
					DType:   []string{"Address"},
				},
				Chain:        chain,
				Name:         alias,
				TokenAddress: txn.ToAddress,
				TxnId:        txn.Hash,
				TxnTime:      txn.Time,
				DType:        []string{"Transaction"},
			}); err != nil {
				s.log.Sugar().Error(err)
				cancel()
				return err
			}
		}
		cancel()
	}
	return nil
}
