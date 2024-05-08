package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostProtectHistoryInfoResponse Response Object
type ListHostProtectHistoryInfoResponse struct {

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 防护状态   - close : 未开启   - opened : 防护中
	ProtectStatus *string `json:"protect_status,omitempty"`

	// total number of static WTPs
	TotalNum *int64 `json:"total_num,omitempty"`

	// data list
	DataList       *[]HostProtectHistoryResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                               `json:"-"`
}

func (o ListHostProtectHistoryInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostProtectHistoryInfoResponse struct{}"
	}

	return strings.Join([]string{"ListHostProtectHistoryInfoResponse", string(data)}, " ")
}
