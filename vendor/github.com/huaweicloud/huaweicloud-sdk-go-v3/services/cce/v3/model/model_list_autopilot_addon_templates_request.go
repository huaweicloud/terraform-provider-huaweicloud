package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotAddonTemplatesRequest Request Object
type ListAutopilotAddonTemplatesRequest struct {

	// 指定的插件名称或插件别名，不填写则查询列表。
	AddonTemplateName *string `json:"addon_template_name,omitempty"`

	// 含义：可接受的最低升级版本  属性：隐藏参数
	BaseUpdateAddonVersion *string `json:"base_update_addon_version,omitempty"`

	// 含义：查询的集群  属性：隐藏参数
	ClusterId *string `json:"cluster_id,omitempty"`

	// 含义：是否获取最新插件  属性：隐藏参数
	Newest *string `json:"newest,omitempty"`

	// 含义：筛选的插件版本  属性：隐藏参数
	Version *string `json:"version,omitempty"`
}

func (o ListAutopilotAddonTemplatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotAddonTemplatesRequest struct{}"
	}

	return strings.Join([]string{"ListAutopilotAddonTemplatesRequest", string(data)}, " ")
}
