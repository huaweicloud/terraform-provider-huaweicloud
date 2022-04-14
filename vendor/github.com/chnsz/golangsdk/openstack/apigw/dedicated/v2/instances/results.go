package instances

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// GetResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

type CreateResp struct {
	Id      string `json:"instance_id"`
	Message string `json:"message"`
}

//Call its Extract method to interpret it as a Instance Id.
func (r CreateResult) Extract() (*CreateResp, error) {
	var s CreateResp
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of a Update operation.
type UpdateResult struct {
	commonResult
}

type Instance struct {
	// Instance ID.
	Id string `json:"id"`
	// Project ID.
	ProjectId string `json:"project_id"`
	// Instance name.
	Name string `json:"instance_name"`
	// Instance status. The value are as following:
	//   Creating, CreateSuccess, CreateFail, Initing, Registering, Running, InitingFailed, RegisterFailed, Installing
	//   InstallFailed, Updating, UpdateFailed, Rollbacking, RollbackSuccess, RollbackFailed, Deleting, DeleteFailed
	//   Unregistering, UnRegisterFailed, CreateTimeout, InitTimeout, RegisterTimeout, InstallTimeout, UpdateTimeout
	//   RollbackTimeout, DeleteTimeout, UnregisterTimeout, Starting, Freezing, Frozen, Restarting, RestartFail
	//   Unhealthy, RestartTimeout
	// The status 'Deleting' is not supported, it's a BUG. --2021/06/15
	Status string `json:"status"`
	// Instance status ID.
	//   1:Creating, 2:CreateSuccess, 3:CreateFail, 4:Initing, 5:Registering, 6:Running, 7:InitingFailed
	//   8:RegisterFailed, 10:Installing, 11:InstallFailed, 12:Updating, 13:UpdateFailed, 20:Rollbacking
	//   21:RollbackSuccess, 22:RollbackFailed, 23:Deleting, 24:DeleteFailed, 25:Unregistering, 26:UnRegisterFailed
	//   27:CreateTimeout, 28:InitTimeout, 29:RegisterTimeout, 30:InstallTimeout, 31:UpdateTimeout
	//   32:RollbackTimeout, 33:DeleteTimeout, 34:UnregisterTimeout, 35:Starting, 36:Freezing, 37:Frozen, 38:Restarting
	//   39:RestartFail, 40:Unhealthy, 41:RestartTimeout
	// Ditto: Issue of status id 23 (Deleting). --2021/06/15
	StatusId int `json:"instance_status"`
	// Instance type.
	Type string `json:"type"`
	// Instance edition.
	Edition string `json:"spec"`
	// Time when the APIG dedicated instance is created, in Unix timestamp format.
	CreateTimestamp int64 `json:"create_time"`
	// Enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Billing mode of the APIG dedicated instance.
	//   0:Pay per use
	//   1:Pay per use
	ChargeMode int `json:"charging_mode"`
	// Yearly/Monthly subscription order ID.
	CbcMetadata string `json:"cbc_metadata"`
	// Description about the APIG dedicated instance.
	Description string `json:"description"`
	// VPC ID.
	VpcId string `json:"vpc_id"`
	// Subnet network ID.
	SubnetId string `json:"subnet_id"`
	// ID of the security group to which the APIG dedicated instance belongs to.
	SecurityGroupId string `json:"security_group_id"`
	// Start time of the maintenance time window in the format "xx:00:00".
	MaintainBegin string `json:"maintain_begin"`
	// End time of the maintenance time window in the format "xx:00:00".
	MaintainEnd string `json:"maintain_end"`
	// VPC ingress private address.
	Ipv4VpcIngressAddress string `json:"ingress_ip"`
	// VPC ingress private address (IPv6).
	Ipv6VpcIngressAddress string `json:"ingress_ip_v6"`
	// ID of the account to which the APIG dedicated instance belongs.
	UserId string `json:"user_id"`
	// EIP bound to the APIG dedicated instance.
	Ipv4IngressEipAddress string `json:"eip_address"`
	// EIP (IPv6).
	Ipv6IngressEipAddress string `json:"eip_ipv6_address"`
	// Public egress address (IPv6).
	Ipv6EgressCidr string `json:"nat_eip_ipv6_cidr"`
	// IP address for public outbound access.
	Ipv4EgressAddress string `json:"nat_eip_address"`
	// Outbound access bandwidth.
	BandwidthSize int `json:"bandwidth_size"`
	// AZs.
	AvailableZoneIds string `json:"available_zone_ids"`
	// Instance version.
	Version string `json:"instance_version"`
	// Supported features.
	SupportedFeatures []string `json:"supported_features"`
}

// Call its Extract method to interpret it as a Instance.
func (r commonResult) Extract() (*Instance, error) {
	var s Instance
	err := r.ExtractInto(&s)
	return &s, err
}

type BaseInstance struct {
	// Instance ID.
	Id string `json:"id"`
	// Project ID
	ProjectId string `json:"project_id"`
	// Instance name.
	Name string `json:"instance_name"`
	// Instance status. The value are as following:
	//   Creating, CreateSuccess, CreateFail, Initing, Registering, Running, InitingFailed, RegisterFailed, Installing
	//   InstallFailed, Updating, UpdateFailed, Rollbacking, RollbackSuccess, RollbackFailed, Deleting, DeleteFailed
	//   Unregistering, UnRegisterFailed, CreateTimeout, InitTimeout, RegisterTimeout, InstallTimeout, UpdateTimeout
	//   RollbackTimeout, DeleteTimeout, UnregisterTimeout, Starting, Freezing, Frozen, Restarting, RestartFail
	//   Unhealthy, RestartTimeout
	// Ditto: Issue of status 'Deleting'. --2021/06/15
	Status string `json:"status"`
	// Instance status ID.
	//   1:Creating, 2:CreateSuccess, 3:CreateFail, 4:Initing, 5:Registering, 6:Running, 7:InitingFailed
	//   8:RegisterFailed, 10:Installing, 11:InstallFailed, 12:Updating, 13:UpdateFailed, 20:Rollbacking
	//   21:RollbackSuccess, 22:RollbackFailed, 23:Deleting, 24:DeleteFailed, 25:Unregistering, 26:UnRegisterFailed
	//   27:CreateTimeout, 28:InitTimeout, 29:RegisterTimeout, 30:InstallTimeout, 31:UpdateTimeout
	//   32:RollbackTimeout, 33:DeleteTimeout, 34:UnregisterTimeout, 35:Starting, 36:Freezing, 37:Frozen, 38:Restarting
	//   39:RestartFail, 40:Unhealthy, 41:RestartTimeout
	// Ditto: Issue of status id 23 (Deleting). --2021/06/15
	StatusId int `json:"instance_status"`
	// Instance type.
	Type string `json:"type"`
	// Instance edition.
	Edition string `json:"spec"`
	// Time when the APIG dedicated instance is created, in Unix timestamp format.
	CreateTimestamp int64 `json:"create_time"`
	// Enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// EIP bound to the APIG dedicated instance.
	Ipv4Address string `json:"eip_address"`
	// Billing mode of the APIG dedicated instance.
	//   0:Pay per use
	//   1:Pay per use
	ChargeMode int `json:"charging_mode"`
	// Yearly/Monthly subscription order ID.
	CbcMetadata string `json:"cbc_metadata"`
}

// InstancePage represents the result of a List operation.
type InstancePage struct {
	pagination.SinglePageBase
}

// Call its Extract method to interpret it as a BaseInstance array.
func ExtractInstances(r pagination.Page) ([]BaseInstance, error) {
	var s []BaseInstance
	err := r.(InstancePage).Result.ExtractIntoSlicePtr(&s, "instances")
	return s, err
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}

// EnableEgressResult represents the result of a EnableEgressAccess operation.
type EnableEgressResult struct {
	golangsdk.Result
}

// UdpateEgressResult represents the result of a UpdateEgressBandwidth operation.
type UdpateEgressResult struct {
	golangsdk.Result
}

type EgressResult struct {
	golangsdk.Result
}

type Egress struct {
	Id               string `json:"id"`
	CloudEipId       string `json:"cloudEipId"`
	CloudEipAddress  string `json:"cloudEipAddress"`
	InstanceId       string `json:"instanceId"`
	CloudBandwidthId string `json:"cloudBandwidthId"`
	BandwidthName    string `json:"bandwidthName"`
	BandwidthSize    int    `json:"bandwidthSize"`
}

// Extract is a method to interpret the response body or json string as an Egress.
func (r UdpateEgressResult) Extract() (*Egress, error) {
	var s Egress
	if r.Err != nil {
		return &s, r.Err
	}
	body, ok := r.Body.(string)
	if ok {
		err := json.Unmarshal([]byte(body), &s)
		return &s, err
	}
	err := r.ExtractInto(&s)
	return &s, err
}

// Extract is a method to interpret the response body as an Egress.
func (r EnableEgressResult) Extract() (*Egress, error) {
	var s Egress
	err := r.ExtractInto(&s)
	return &s, err
}

// DisableEgressResult represents the result of a DisableEgressAccess operation.
type DisableEgressResult struct {
	golangsdk.ErrResult
}

// EnableIngressResult represents the result of a EnableIngressAccess operation.
type EnableIngressResult struct {
	commonResult
}

type Ingress struct {
	Id          string `json:"eip_id"`
	EipAddress  string `json:"eip_address"`
	Status      string `json:"eip_status"`
	Ipv6Address string `json:"eip_ipv6_address"`
}

// Call its Extract method to interpret it as a Ingress.
func (r EnableIngressResult) Extract() (*Ingress, error) {
	var s Ingress
	err := r.ExtractInto(&s)
	return &s, err
}

// DisableIngressResult represents the result of a DisableIngressAccess operation.
type DisableIngressResult struct {
	golangsdk.ErrResult
}
