package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListIsolatedFileResponse Response Object
type ListIsolatedFileResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 已隔离文件详情
	DataList       *[]IsolatedFileResponseInfo `json:"data_list,omitempty"`
	HttpStatusCode int                         `json:"-"`
}

func (o ListIsolatedFileResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListIsolatedFileResponse struct{}"
	}

	return strings.Join([]string{"ListIsolatedFileResponse", string(data)}, " ")
}
