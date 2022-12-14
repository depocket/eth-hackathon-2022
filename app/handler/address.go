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

	c.JSON(200, TransformFlowResponse(flow, jsonParams.Address))
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

	c.JSON(200, TransformFlowResponse(flow, jsonParams.Address))
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
	c.JSON(200, TransformFlowResponse(flow, jsonParams.Address))

}

func (h *AddressHandler) Path(c *gin.Context) {
	var jsonParams model.PathRequest
	if err := c.ShouldBindJSON(&jsonParams); err != nil {
		utils.Response(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), utils.GeneralTimeout)
	defer cancel()

	path, err := h.service.Path(ctx, jsonParams)
	if err != nil {
		utils.Response(c, err)
		return
	}

	c.JSON(200, TransformPathResponse(path, jsonParams.ToAddress, jsonParams.FromAddress))

}

func TransformFlowResponse(flow *model.ResponseFlow, address string) model.FlowTransformed {
	nodes := make(map[string]model.Node, 0)
	edges := make(map[string]model.Edge, 0)
	resNode := make([]model.Node, 0)
	resEdge := make([]model.Edge, 0)
	for _, data := range flow.Data {
		identifyFlow(data, nodes, edges, utils.ColorMainNode)
	}
	for _, v := range nodes {
		if v.Label == address {
			v.Color = utils.ColorMainNode
		}
		resNode = append(resNode, v)
	}
	for _, v := range edges {
		resEdge = append(resEdge, v)
	}
	return model.FlowTransformed{
		Nodes: resNode,
		Edges: resEdge,
	}
}

func TransformPathResponse(path *model.ResponsePath, addressFrom string, addrerssTo string) model.PathTransformed {
	trans := make(map[string]string, 0)
	resNode := make([]model.Node, 0)
	for _, n := range path.Node {
		if n.Address == "" {
			trans[n.UID] = n.Name
		} else {
			color := utils.ColorMainNode
			if n.Address == addressFrom || n.Address == addrerssTo {
				color = utils.ColorFromToNode
			}
			resNode = append(resNode, model.Node{
				Id:    n.UID,
				Label: n.Address,
				Title: "address",
				Color: color,
			})
		}
	}

	//-------------------------------------
	resEdge := make([]model.Edge, 0)
	for _, p := range path.Path {
		edge := make(map[int]string, 0)
		identifyPath(p, edge, utils.ToPointer(0))
		for i := 0; i < len(edge)-2; i += 2 {
			resEdge = append(resEdge, model.Edge{
				ID:    edge[i+1],
				From:  edge[i],
				To:    edge[i+2],
				Label: trans[edge[i+1]],
				Smooth: model.SmoothEdge{
					Type:      utils.SmoothType(),
					Roundness: utils.SmoothRoundness(),
				},
			})
		}
	}
	return model.PathTransformed{
		Nodes: resNode,
		Edges: resEdge,
	}
}

func identifyFlow(output model.AddressDgraphResponse, nodes map[string]model.Node, edges map[string]model.Edge, color string) {
	if output.UID != "" {
		nodes[output.UID] = model.Node{
			Id:    output.UID,
			Label: output.Address,
			Title: "address",
			Color: color,
		}
		for _, sender := range output.Sender {
			if sender.UID != "" {
				if output.UID != "" && sender.Recipient.UID != "" {
					edges[output.UID+sender.Recipient.UID+sender.Name] = model.Edge{
						ID:       sender.UID,
						From:     output.UID,
						To:       sender.Recipient.UID,
						Label:    sender.Name,
						Relation: "out",
						Smooth: model.SmoothEdge{
							Type:      utils.SmoothType(),
							Roundness: utils.SmoothRoundness(),
						},
					}
					color = utils.ColorSender
					identifyFlow(sender.Sender, nodes, edges, color)
					identifyFlow(sender.Recipient, nodes, edges, color)
				}
			}
		}
		for _, recipient := range output.Recipient {
			if recipient.UID != "" {
				if recipient.Sender.UID != "" && output.UID != "" {
					edges[recipient.Sender.UID+output.UID+recipient.Name] = model.Edge{
						ID:       recipient.UID,
						From:     recipient.Sender.UID,
						To:       output.UID,
						Label:    recipient.Name,
						Relation: "in",
						Smooth: model.SmoothEdge{
							Type:      utils.SmoothType(),
							Roundness: utils.SmoothRoundness(),
						},
					}
					color = utils.ColorRecipient
					identifyFlow(recipient.Sender, nodes, edges, color)
					identifyFlow(recipient.Recipient, nodes, edges, color)
				}
			}
		}
	}
}
func identifyPath(path map[string]interface{}, edge map[int]string, count *int) {
	edge[*count] = path["uid"].(string)
	*count++
	if v, ok := path["~recipient"]; ok {
		identifyPath(v.(map[string]interface{}), edge, count)
	} else if v, ok := path["recipient"]; ok {
		identifyPath(v.(map[string]interface{}), edge, count)
	} else if v, ok := path["~sender"]; ok {
		identifyPath(v.(map[string]interface{}), edge, count)
	} else if v, ok := path["sender"]; ok {
		identifyPath(v.(map[string]interface{}), edge, count)
	} else {
		return
	}
}
