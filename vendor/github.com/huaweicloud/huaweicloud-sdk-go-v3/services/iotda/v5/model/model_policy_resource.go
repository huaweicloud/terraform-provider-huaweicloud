package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PolicyResource 预调配模板设备策略资源详情结构体。
type PolicyResource struct {

	// **参数说明**：设备需要绑定的策略id列表
	PolicyIds *[]string `json:"policy_ids,omitempty"`
}

func (o PolicyResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyResource struct{}"
	}

	return strings.Join([]string{"PolicyResource", string(data)}, " ")
}
