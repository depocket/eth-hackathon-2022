package handler

import (
	"context"
	"depocket.io/app/model"
	"depocket.io/app/service"
	"depocket.io/app/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AddressHandler struct {
	log     *zap.Logger
	service service.AddressInterface
}

func NewAddressHandler(r *gin.RouterGroup, log *zap.Logger, service service.AddressInterface) {
	h := &AddressHandler{
		log:     log,
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

	ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
	defer cancel()
	flow, err := h.service.FullFlow(ctx, jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}

	c.JSON(200, TransformFlowResponse(flow))
}

func (h *AddressHandler) InFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
	defer cancel()
	flow, err := h.service.InFlow(ctx, jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}

	c.JSON(200, TransformFlowResponse(flow))
}

func (h *AddressHandler) OutFlow(c *gin.Context) {
	var jsonParams model.FlowRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
	defer cancel()
	flow, err := h.service.OutFlow(ctx, jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, TransformFlowResponse(flow))

}

func (h *AddressHandler) Path(c *gin.Context) {
	var jsonParams model.PathRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
	defer cancel()
	recommend, err := h.service.Path(ctx, jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, recommend)

}

func TransformFlowResponse(flow *model.ResponseFlow) model.FlowTransformed {
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
	return model.FlowTransformed{
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
				if sender.Sender.UID != "" && sender.Recipient.UID != "" {
					edges[sender.Sender.UID+sender.Recipient.UID+sender.Name] = model.Edge{
						ID:       sender.TxnTime.Unix(),
						From:     sender.Sender.UID,
						To:       sender.Recipient.UID,
						Label:    sender.Name,
						Relation: "out",
					}
					identify(sender.Sender, nodes, edges)
					identify(sender.Recipient, nodes, edges)
				}
			}
		}
		for _, recipient := range output.Recipient {
			if recipient.UID != "" {
				if recipient.Sender.UID != "" && recipient.Recipient.UID != "" {
					edges[recipient.Sender.UID+recipient.Recipient.UID+recipient.Name] = model.Edge{
						ID:       recipient.TxnTime.Unix(),
						From:     recipient.Sender.UID,
						To:       recipient.Recipient.UID,
						Label:    recipient.Name,
						Relation: "in",
					}
					identify(recipient.Sender, nodes, edges)
					identify(recipient.Recipient, nodes, edges)
				}
			}
		}
	}
}
