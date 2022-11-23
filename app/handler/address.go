package handler

import (
	"depocket.io/app/model"
	"depocket.io/app/service"
	"depocket.io/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddressHandler struct {
	log     *zap.Logger
	db      *gorm.DB
	service service.AddressInterface
}

func NewAddressHandler(r *gin.RouterGroup, log *zap.Logger, db *gorm.DB, service service.AddressInterface) {
	h := &AddressHandler{
		log:     log,
		db:      db,
		service: service,
	}
	ar := r.Group("/address")
	ar.POST("/full-flow", h.FullFlow)
	ar.POST("/in-flow", h.InFlow)
	ar.POST("/out-flow", h.OutFlow)
	ar.POST("/path", h.Path)
}

func (h *AddressHandler) FullFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	flow, err := h.service.FullFlow(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, flow)
}

func (h *AddressHandler) InFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	recommend, err := h.service.InFlow(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, recommend)
}

func (h *AddressHandler) OutFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	recommend, err := h.service.OutFlow(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, recommend)

}

func (h *AddressHandler) Path(c *gin.Context) {
	var jsonParams model.PathRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	recommend, err := h.service.Path(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, recommend)

}
