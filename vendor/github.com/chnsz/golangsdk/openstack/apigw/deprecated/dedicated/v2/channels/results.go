package channels

import (
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

// UdpateResult represents the result of a update operation.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a update operation.
type GetResult struct {
	commonResult
}

type VpcChannel struct {
	// VPC channel name.
	// A VPC channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// VPC channel type.
	// 1: private network ELB channel (to be deprecated)
	// 2: fast channel with the load balancing function
	Type int `json:"type"`
	// Host port of the VPC channel.
	// This parameter is valid only when the VPC channel type is set to 2.
	// This parameter is required if the VPC channel type is set to 2.
	// The value range is 1–65535.
	Port int `json:"port"`
	// Distribution algorithm.
	// 1: WRR (default)
	// 2: WLC
	// 3: SH
	// 4: URI hashing
	// This parameter is mandatory if the VPC channel type is set to 2.
	BalanceStrategy int `json:"balance_strategy"`
	// Member type of the VPC channel.
	// ip
	// ecs (default)
	// This parameter is required if the VPC channel type is set to 2.
	MemberType string `json:"member_type"`
	// Time when the VPC channel is created.
	CreateTime string `json:"create_time"`
	// VPC channel ID.
	Id string `json:"id"`
	// VPC channel status.
	// 1: normal
	// 2: abnormal
	Status int `json:"status"`
	// ID of a private network ELB channel.
	// This parameter is valid only when the VPC channel type is set to 1.
	ElbId string `json:"elb_id"`
	// Backend server list. Only one backend server is included if the VPC channel type is set to 1.
	Members []MemberInfo `json:"members"`
	// Health check details.
	VpcHealthConfig VpcHealthConfig `json:"vpc_health_config"`
}

func (r commonResult) Extract() (*VpcChannel, error) {
	var s VpcChannel
	err := r.ExtractInto(&s)
	return &s, err
}

// The ChannelPage represents the result of a List operation.
type ChannelPage struct {
	pagination.SinglePageBase
}

// ExtractChannels its Extract method to interpret it as a channel array.
func ExtractChannels(r pagination.Page) ([]VpcChannel, error) {
	var s []VpcChannel
	err := r.(ChannelPage).Result.ExtractIntoSlicePtr(&s, "vpc_channels")
	return s, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type MemberResult struct {
	golangsdk.Result
}

type AddBackendResult struct {
	MemberResult
}

type GetBackendResult struct {
	MemberResult
}

type Member struct {
	// VPC channel name.
	// A VPC channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// VPC channel type.
	// 1: private network ELB channel (to be deprecated)
	// 2: fast channel with the load balancing function
	Type int `json:"type"`
	// Host port of the VPC channel.
	// This parameter is valid only when the VPC channel type is set to 2. The value range is 1–65535.
	// This parameter is required if the VPC channel type is set to 2.
	Port int `json:"port"`
	// Distribution algorithm.
	// 1: WRR (default)
	// 2: WLC
	// 3: SH
	// 4: URI hashing
	// This parameter is mandatory if the VPC channel type is set to 2.
	BalanceStrategy int `json:"balance_strategy"`
	// Member type of the VPC channel.
	// ip
	// ecs (default)
	// This parameter is required if the VPC channel type is set to 2.
	MemberType string `json:"member_type"`
	// Time when the VPC channel is created.
	CreateTime string `json:"create_time"`
	// VPC channel ID.
	Id string `json:"id"`
	// VPC channel status.
	// 1: normal
	// 2: abnormal
	Status int `json:"status"`
	// ID of a private network ELB channel.
	// This parameter is valid only when the VPC channel type is set to 1.
	ElbId string `json:"elb_id"`
}

func (r MemberResult) Extract() ([]Member, error) {
	var s []Member
	err := r.ExtractIntoStructPtr(&s, "members")
	return s, err
}

type RemoveResult struct {
	golangsdk.ErrResult
}
