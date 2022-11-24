package handler

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/service"
	"depocket.io/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SyncAddressHandler struct {
	log     *zap.Logger
	service service.SyncAddressInterface
}

func NewSyncAddressHandler(r *gin.RouterGroup, log *zap.Logger, service service.SyncAddressInterface) {
	h := &SyncAddressHandler{
		log:     log,
		service: service,
	}
	ar := r.Group("/dgraph")
	ar.GET("/sync", h.Sync)
}

func (h *SyncAddressHandler) Sync(c *gin.Context) {
	chain := c.Query("chain")
	limitRq := c.Query("limit")
	limit, err := strconv.Atoi(limitRq)
	if err != nil {
		limit = 1000
	}
	if chain == "" {
		utils.Response(c, utils.Error{
			Status:  http.StatusBadRequest,
			Message: "Empty chain request",
		})
		return
	}
	h.FetchTransactions(chain, limit)
	c.JSON(http.StatusOK, fmt.Sprintf("Syncing transactions from %s", chain))
}

func (h *SyncAddressHandler) FetchTransactions(chain string, limit int) {
	s := time.Now()
	fetching := true
	cursor := ""
	totalFetch := 0
	retry := 0
	fetchRequest := model.TransactionRequest{
		Chain:         utils.ToPointer(chain),
		Limit:         utils.ToPointer(100),
		DecodedAction: utils.ToPointer("transfer,transferFrom"),
		Decoded:       utils.ToPointer(true),
	}
	h.log.Sugar().Info(chain + " fetching...")
	for fetching && totalFetch <= limit {
		// get trans
		t := time.Now()
		fetchRequest.Cursor = utils.ToPointer(cursor)
		txns, err := utils.FetchDepocketTransaction(fetchRequest)
		if err != nil {
			h.log.Sugar().Warn(err)
			if retry < 3 {
				retry++
				continue
			}
			h.log.Sugar().Error("Too many request/retry ", err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
		// get tokens in all trans
		addrs := make(map[string]string, 0)
		tokenAddrs := make([]string, 0)
		for _, txn := range txns.Data {
			if _, ok := addrs[txn.ToAddress]; !ok {
				tokenAddrs = append(tokenAddrs, txn.ToAddress)
			}
			addrs[txn.ToAddress] = txn.ToAddress
		}

		tokens, err := utils.FetchDepocketToken(model.TokenRequest{
			Chain:     utils.ToPointer(chain),
			Addresses: utils.ToPointer(strings.Join(tokenAddrs, ",")),
		})
		if err != nil {
			h.log.Sugar().Warn(err)
			cancel()
			return
		}
		tokenSymbols := make(map[string]model.Token, 0)
		for _, token := range tokens.Data {
			tokenSymbols[token.Address] = token
		}

		if err := h.service.SyncAddress(ctx, chain, txns.Data, tokenSymbols); err != nil {
			h.log.Sugar().Warn(err)
			cancel()
			return
		}
		fetching = txns.HasNextPage
		cursor = txns.NextPageCursor
		totalFetch += len(txns.Data)
		h.log.Sugar().Infof("Synced %v transaction to dgraph in %v", totalFetch, time.Now().Sub(t))
		cancel()
	}
	h.log.Sugar().Infof("%s fetched after %v", chain, time.Now().Sub(s))
}
