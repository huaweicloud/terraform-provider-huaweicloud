package flavors

type NodeTypes struct {
	NodeTypes []NodeType `json:"node_types"`
}

type NodeType struct {
	Detail   []TypeDetail `json:"detail"`
	Id       string       `json:"id"`
	SpecName string       `json:"spec_name"`
}

type TypeDetail struct {
	Unit  string `json:"unit"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
