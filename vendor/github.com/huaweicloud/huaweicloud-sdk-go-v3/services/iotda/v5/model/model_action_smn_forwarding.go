package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 发送给SMN消息结构
type ActionSmnForwarding struct {

	// **参数说明**：SMN服务对应的region区域。
	RegionName string `json:"region_name"`

	// **参数说明**：SMN服务对应的projectId信息。
	ProjectId string `json:"project_id"`

	// **参数说明**：SMN服务对应的主题名称。
	ThemeName string `json:"theme_name"`

	// **参数说明**：SMN服务对应的topic的主题URN。
	TopicUrn string `json:"topic_urn"`

	// **参数说明**：短信或邮件的内容。。
	MessageContent string `json:"message_content"`

	// **参数说明**：短信或邮件的主题。。
	MessageTitle string `json:"message_title"`
}

func (o ActionSmnForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ActionSmnForwarding struct{}"
	}

	return strings.Join([]string{"ActionSmnForwarding", string(data)}, " ")
}
