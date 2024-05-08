package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppsResponse Response Object
type ListAppsResponse struct {

	// 软件总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 软件列表
	DataList       *[]AppResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListAppsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppsResponse struct{}"
	}

	return strings.Join([]string{"ListAppsResponse", string(data)}, " ")
}
