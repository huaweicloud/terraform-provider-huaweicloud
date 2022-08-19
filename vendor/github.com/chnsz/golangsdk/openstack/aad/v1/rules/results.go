package rules

// BatchResp is the structure that represents the API request response for the create and delete methods.
type BatchResp struct {
	SuccessNum int `json:"successNum"`
}

// ListResp is the structure that represents the API request response for list method.
type ListResp struct {
	Count int    `json:"count"`
	Rules []Rule `json:"rules"`
}

// Rule is the structure that represents the details of the forward rule.
type Rule struct {
	ID              string `json:"rule_id"`
	ForwardProtocol string `json:"forward_protocol"`
	ForwardPort     int    `json:"forward_port"`
	SourcePort      int    `json:"source_port"`
	SourceIp        string `json:"source_ip"`
	LbMethod        string `json:"lb_method"`
	Status          int    `json:"status"`
}
