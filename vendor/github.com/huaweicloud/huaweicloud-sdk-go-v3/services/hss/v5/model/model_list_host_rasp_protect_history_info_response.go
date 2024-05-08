package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostRaspProtectHistoryInfoResponse Response Object
type ListHostRaspProtectHistoryInfoResponse struct {

	// total number of dynamic WTPs
	TotalNum *int64 `json:"total_num,omitempty"`

	// data list
	DataList       *[]HostRaspProtectHistoryResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                                   `json:"-"`
}

func (o ListHostRaspProtectHistoryInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostRaspProtectHistoryInfoResponse struct{}"
	}

	return strings.Join([]string{"ListHostRaspProtectHistoryInfoResponse", string(data)}, " ")
}
