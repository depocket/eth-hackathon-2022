package model

type FlowTransformed struct {
	Data  []AddressDgraphResponse `json:"data"`
	Nodes []Node                  `json:"nodes"`
	Edges []Edge                  `json:"edges"`
}

type PathTransformed struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Title string `json:"title"`
}

type Edge struct {
	ID       string `json:"id"`
	From     string `json:"from"`
	To       string `json:"to"`
	Label    string `json:"label"`
	Relation string `json:"relation"`
}
