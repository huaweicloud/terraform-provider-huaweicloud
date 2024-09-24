package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowInstanceMessagesRequest Request Object
type ShowInstanceMessagesRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// Topic名称。  Topic名称必须以字母开头且只支持大小写字母、中横线、下划线以及数字。
	Topic string `json:"topic"`

	// 是否按照时间排序。
	Asc *bool `json:"asc,omitempty"`

	// 开始时间。  Unix毫秒时间戳。  查询消息偏移量时，为必选参数。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。  Unix毫秒时间戳。  查询消息偏移量时，为必选参数。
	EndTime *string `json:"end_time,omitempty"`

	// 每一页显示的message数量。
	Limit *string `json:"limit,omitempty"`

	// 页数。
	Offset *string `json:"offset,omitempty"`

	// 是否下载。
	Download *bool `json:"download,omitempty"`

	// 消息偏移量。  **查询消息内容时，为必选参数。**  若start_time、end_time参数不为空，该参数无效。
	MessageOffset *string `json:"message_offset,omitempty"`

	// 分区。  **查询消息内容时，为必选参数。**  若start_time、end_time参数不为空，该参数无效。
	Partition *string `json:"partition,omitempty"`

	// 关键词。 取值范围为0~50。
	Keyword *string `json:"keyword,omitempty"`
}

func (o ShowInstanceMessagesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceMessagesRequest struct{}"
	}

	return strings.Join([]string{"ShowInstanceMessagesRequest", string(data)}, " ")
}
