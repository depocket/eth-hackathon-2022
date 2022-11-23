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

	c.JSON(200, TransformResponse(flow))
}

func (h *AddressHandler) InFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	flow, err := h.service.InFlow(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}

	c.JSON(200, TransformResponse(flow))
}

func (h *AddressHandler) OutFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	flow, err := h.service.OutFlow(jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, TransformResponse(flow))

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

func TransformResponse(flow *model.ResponseFlow) model.ResponseTransformed {
	nodes := make(map[string]model.Node, 0)
	edges := make(map[string]model.Edge, 0)
	resNode := make([]model.Node, 0)
	resEdge := make([]model.Edge, 0)
	for _, data := range flow.Data {
		identify(data, nodes, edges)
	}
	for _, v := range nodes {
		resNode = append(resNode, v)
	}
	for _, v := range edges {
		resEdge = append(resEdge, v)
	}
	return model.ResponseTransformed{
		Data:  flow.Data,
		Nodes: resNode,
		Edges: resEdge,
	}
}

func identify(output model.AddressDgraphResponse, nodes map[string]model.Node, edges map[string]model.Edge) {
	if output.UID != "" {
		nodes[output.UID] = model.Node{
			Id:    output.UID,
			Label: output.Address,
			Title: "address",
		}
		for _, sender := range output.Sender {
			if sender.UID != "" {
				nodes[sender.UID] = model.Node{
					Id:    sender.UID,
					Label: sender.Name,
					Title: "transactions",
				}
				edges[output.UID+sender.UID] = model.Edge{
					From: output.UID,
					To:   sender.UID,
				}
				if sender.Sender.UID != "" {
					edges[sender.UID+sender.Sender.UID] = model.Edge{
						From: sender.UID,
						To:   sender.Sender.UID,
					}
					identify(sender.Sender, nodes, edges)
				}
				if sender.Recipient.UID != "" {
					edges[sender.UID+sender.Recipient.UID] = model.Edge{
						From: sender.UID,
						To:   sender.Recipient.UID,
					}
					identify(sender.Recipient, nodes, edges)
				}
			}
		}
		for _, recipient := range output.Recipient {
			if recipient.UID != "" {
				nodes[recipient.UID] = model.Node{
					Id:    recipient.UID,
					Label: recipient.Name,
					Title: "transactions",
				}
				edges[output.UID+"-"+recipient.UID] = model.Edge{
					From: output.UID,
					To:   recipient.UID,
				}
				if recipient.Sender.UID != "" {
					edges[recipient.UID+recipient.Sender.UID] = model.Edge{
						From: recipient.UID,
						To:   recipient.Sender.UID,
					}
					identify(recipient.Sender, nodes, edges)
				}
				if recipient.Recipient.UID != "" {
					edges[recipient.UID+recipient.Recipient.UID] = model.Edge{
						From: recipient.UID,
						To:   recipient.Recipient.UID,
					}
					identify(recipient.Recipient, nodes, edges)
				}
			}
		}
	}
}
