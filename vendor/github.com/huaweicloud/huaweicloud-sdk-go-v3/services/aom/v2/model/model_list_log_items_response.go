package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListLogItemsResponse struct {

	// 响应码,SVCSTG_AMS_2000000代表正常返回。
	ErrorCode *string `json:"errorCode,omitempty"`

	// 响应信息描述。
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// 查询结果元数据信息，包括返回总数及结果。
	Result         *string `json:"result,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListLogItemsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogItemsResponse struct{}"
	}

	return strings.Join([]string{"ListLogItemsResponse", string(data)}, " ")
}
