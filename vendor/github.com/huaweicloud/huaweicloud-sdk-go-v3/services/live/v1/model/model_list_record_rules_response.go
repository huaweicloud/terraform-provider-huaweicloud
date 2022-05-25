package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListRecordRulesResponse struct {

	// 查询结果的总元素数量
	Total *int32 `json:"total,omitempty"`

	// 录制配置数组
	RecordConfig   *[]RecordRule `json:"record_config,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListRecordRulesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecordRulesResponse struct{}"
	}

	return strings.Join([]string{"ListRecordRulesResponse", string(data)}, " ")
}
