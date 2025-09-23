package channels

import (
	"github.com/chnsz/golangsdk/pagination"
)

// Channel is the structure that represents the channel detail.
type Channel struct {
	// Channel name.
	// A channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// Host port of the channel.
	// The value range is 1–65535.
	Port int `json:"port"`
	// Distribution algorithm.
	// 1: WRR (default)
	// 2: WLC
	// 3: SH
	// 4: URI hashing
	BalanceStrategy int `json:"balance_strategy"`
	// Member type of the channel.
	// ip
	// ecs (default)
	MemberType string `json:"member_type"`
	// Channel type.
	// + 2: Server type.
	// + 3: Microservice type.
	Type int `json:"type"`
	// builtin: server type
	// + microservice: microservice type
	// + reference: reference load balance
	VpcChannelType string `json:"vpc_channel_type"`
	// Dictionary code of the channel.
	// The value can contain letters, digits, hyphens (-), underscores (_), and periods (.).
	DictCode string `json:"dict_code"`
	// Time when the channel is created.
	CreateTime string `json:"create_time"`
	// Channel ID.
	ID string `json:"id"`
	// Channel status.
	// 1: normal
	// 2: abnormal
	Status int `json:"status"`
	// Backend server groups of the channel.
	MemberGroups []MemberGroup `json:"member_groups"`
	// Backend server list. Only one backend server is included if the channel type is set to 1.
	Members []MemberInfo `json:"members"`
	// Health check details.
	VpcHealthConfig *VpcHealthConfig `json:"vpc_health_config"`
	// Microservice details.
	MicroserviceConfig *MicroserviceConfig `json:"microservice_info"`
}

// ChannelPage is a single page maximum result representing a query by offset page.
type ChannelPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ChannelPage struct is empty.
func (b ChannelPage) IsEmpty() (bool, error) {
	arr, err := ExtractChannels(b)
	return len(arr) == 0, err
}

// ExtractChannels is a method to extract the list of channels.
func ExtractChannels(r pagination.Page) ([]Channel, error) {
	var s []Channel
	err := r.(ChannelPage).Result.ExtractIntoSlicePtr(&s, "vpc_channels")
	return s, err
}

type Member struct {
	// Channel name.
	// A channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// Channel type.
	// 1: private network ELB channel (to be deprecated)
	// 2: fast channel with the load balancing function
	Type int `json:"type"`
	// Host port of the channel.
	// This parameter is valid only when the channel type is set to 2. The value range is 1–65535.
	// This parameter is required if the channel type is set to 2.
	Port int `json:"port"`
	// Distribution algorithm.
	// 1: WRR (default)
	// 2: WLC
	// 3: SH
	// 4: URI hashing
	// This parameter is mandatory if the channel type is set to 2.
	BalanceStrategy int `json:"balance_strategy"`
	// Member type of the channel.
	// ip
	// ecs (default)
	// This parameter is required if the channel type is set to 2.
	MemberType string `json:"member_type"`
	// Time when the channel is created.
	CreateTime string `json:"create_time"`
	// Channel ID.
	Id string `json:"id"`
	// Channel status.
	// 1: normal
	// 2: abnormal
	Status int `json:"status"`
	// ID of a private network ELB channel.
	// This parameter is valid only when the channel type is set to 1.
	ElbId string `json:"elb_id"`
}

// MemberPage is a single page maximum result representing a query by offset page.
type MemberPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a MemberPage struct is empty.
func (b MemberPage) IsEmpty() (bool, error) {
	arr, err := ExtractMembers(b)
	return len(arr) == 0, err
}

// ExtractMembers is a method to extract the list of backend members.
func ExtractMembers(r pagination.Page) ([]MemberInfo, error) {
	var s []MemberInfo
	err := r.(MemberPage).Result.ExtractIntoSlicePtr(&s, "members")
	return s, err
}
