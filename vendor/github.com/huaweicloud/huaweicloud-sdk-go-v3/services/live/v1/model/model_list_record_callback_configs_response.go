package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRecordCallbackConfigsResponse struct {

	// 查询结果的总元素数量
	Total *int32 `json:"total,omitempty"`

	// 回调配置
	CallbackConfig *[]RecordCallbackConfig `json:"callback_config,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListRecordCallbackConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecordCallbackConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListRecordCallbackConfigsResponse", string(data)}, " ")
}
