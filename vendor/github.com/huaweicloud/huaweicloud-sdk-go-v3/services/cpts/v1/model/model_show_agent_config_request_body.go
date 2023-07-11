package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAgentConfigRequestBody 获取探针配置信息请求体
type ShowAgentConfigRequestBody struct {

	// 应用id
	AppId int32 `json:"app_id"`

	// 探针的内网地址
	Address string `json:"address"`

	// 探针的版本
	Version string `json:"version"`

	// 探针id，非必填，不填是注册探针，填了是更新探针配置
	AgentId *string `json:"agent_id,omitempty"`
}

func (o ShowAgentConfigRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAgentConfigRequestBody struct{}"
	}

	return strings.Join([]string{"ShowAgentConfigRequestBody", string(data)}, " ")
}
