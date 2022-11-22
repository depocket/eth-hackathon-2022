package handler

import (
	"depocket.io/app/model"
	"depocket.io/app/service"
	"depocket.io/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type SyncAddressHandler struct {
	log     *zap.Logger
	db      *gorm.DB
	service service.SyncAddressInterface
}

func NewSyncAddressHandler(r *gin.RouterGroup, log *zap.Logger, db *gorm.DB, service service.SyncAddressInterface) {
	h := &SyncAddressHandler{
		log:     log,
		db:      db,
		service: service,
	}
	ar := r.Group("/dgraph")
	ar.GET("/sync", h.Sync)
}

func (h *SyncAddressHandler) Sync(c *gin.Context) {

	// mock transactions
	var txns []model.Transaction
	if err := h.db.Model(&model.Transaction{}).
		Select("from_address, to_address, hash, decoded_input, time").
		Where("action = ? AND decoded_input IS NOT NULL", utils.TransferAction).
		Limit(10000).
		Order("block_number DESC").
		Find(&txns).Error; err != nil {
		h.log.Sugar().Error(err)
		c.Error(err)
		return
	}

	var addrs []string
	for _, txn := range txns {
		addrs = append(addrs, txn.ToAddress)
	}
	var tokens []model.Token
	if err := h.db.Model(&model.Token{}).
		Select("symbol, address, decimals").
		Where("address IN (?)", addrs).
		Find(&tokens).Error; err != nil {
		h.log.Sugar().Error(err)
		c.Error(err)
		return
	}
	tokenSymbols := make(map[string]model.Token, 0)
	for _, token := range tokens {
		tokenSymbols[token.Address] = token
	}
	if err := h.service.SyncAddress(c, txns, tokenSymbols); err != nil {
		h.log.Sugar().Error(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, len(txns))
}
