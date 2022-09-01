package resources

// QueryResp is the structure that represents the API response of 'Get' method.
type QueryResp struct {
	// Error code.
	ErrorCode string `json:"error_code"`
	// Error message.
	ErrorMsg string `json:"error_msg"`
	// List of the prepaid resources.
	Resources []Resource `json:"data"`
	// The total number of resources queried.
	TotalCount int `json:"total_count"`
}

// Resource is the structure that represents the detail of the prepaid resource.
type Resource struct {
	// The internal ID of the resource to be activated. The ID generated after the resource is activated is resource_id.
	ID string `json:"id"`
	// Resource ID.
	ResourceId string `json:"resource_id"`
	// Resource name.
	ResourceName string `json:"resource_name"`
	// Cloud service area code, for example: "cn-north-1".
	Region string `json:"region_code"`
	// Cloud service type code, for example: the cloud service type code of OBS is "hws.service.type.obs".
	ServiceTypeCode string `json:"service_type_code"`
	// Resource type name, for example: the VM of ECS is "hws.resource.type.vm".
	ServieTypeName string `json:"service_type_name"`
	// Resource specifications for cloud service. If it is a resource specification of a VM, the specification will
	// contain ".win" or ".linux", such as "s2.small.1.linux".
	ResourceTypeCode string `json:"resource_type_code"`
	// Cloud service type name, for example: the name of the cloud service type of ECS is "弹性云服务器".
	ResourceTypeName string `json:"resource_type_name"`
	// Resource type name, for example: the resource type name of ECS is "云主机".
	ResourceSpecCode string `json:"resource_spec_code"`
	// Project ID.
	ProjectId string `json:"project_id"`
	// Product ID.
	ProductId string `json:"product_id"`
	// Parent resource ID.
	ParentResourceId string `json:"parent_resource_id"`
	// Whether it is the main resource.
	// + 0: non-primary resource
	// + 1: Main resource
	IsMainResource int `json:"is_main_resource"`
	// Status code.
	// + 2: in use
	// + 3: Closed (the page does not display this state)
	// + 4: Frozen
	// + 5: Expired
	Status int `json:"status"`
	// Resource effective time.
	EffectiveTime string `json:"effective_time"`
	// Resource expire time.
	ExpireTime string `json:"expire_time"`
	// Resource expire policy.
	// + 0: Expires into the grace period
	// + 1: Expired to on-demand
	// + 2: Automatic deletion after expiration (direct deletion from effective)
	// + 3: Automatic renewal after expiration
	// + 4: Freeze after expiration
	// + 5: Delete after expiration (delete from retention period)
	ExpirePolicy int `json:"expire_policy"`
}
