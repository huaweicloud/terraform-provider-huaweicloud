package clouds

// createResp represents the result of the 'Create' function.
type createResp struct {
	// Order ID.
	OrderId string `json:"orderId"`
}

// Instance is the structure that represents the details of the cloud WAF instance.
type Instance struct {
	// The type of the cloud WAF.
	// + -2: Freezen
	// + -1: Not subscribed
	// + 2: Standard edition
	// + 3: Professional edition
	// + 4: Platinum edition
	// + 7: Introduction edition.
	// + 22: Post paid edition.
	Type int `json:"type"`
	// The list of all cloud WAF resources.
	Resources []Resource `json:"resources"`
	// Whether the current user is new user.
	IsNewUser bool `json:"isNewUser"`
	// The subscribe information of the dedicated mode.
	Premium Premium `json:"premium"`
}

// Resource is the structure that represents the details of the WAF resource.
type Resource struct {
	// Resource ID.
	ID string `json:"resourceId"`
	// Resource type.
	// + hws.resource.type.waf
	// + hws.resource.type.waf.domain
	// + hws.resource.type.waf.bandwidth
	// + hws.resource.type.waf.rule
	Type string `json:"resourceType"`
	// The number of the resources.
	Size int `json:"resourceSize"`
	// The cloud service type corresponding to the cloud service product.
	CloudServiceType string `json:"cloudServiceType"`
	// The specification code of the resource.
	SpecCode string `json:"resourceSpecCode"`
	// The current status of the resource.
	Status int `json:"status"`
	// The expire time of the resource.
	// + 0: normal
	// + 1: freezen
	ExpireTime string `json:"expireTime"`
}

// Premium is the structure that represents the configuration of the dedicated mode.
type Premium struct {
	// Whether the dedicated mode is opened.
	Purchased bool `json:"purchased"`
	// The number of the dedicated instances.
	Total int `json:"total"`
	// The number of the ELB instances.
	ELB int `json:"elb"`
	// The number of the dedicated WAFs.
	Dedicated int `json:"dedicated"`
}

// updateResp represents the result of the 'Update' function.
type updateResp struct {
	// Order ID.
	OrderId string `json:"orderId"`
}
