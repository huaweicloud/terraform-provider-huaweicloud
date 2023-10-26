package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EsHealthIpgroupResource struct {

	// 监听器关联的访问控制组的id。
	IpgroupId *string `json:"ipgroup_id,omitempty"`

	// 访问控制组的状态。
	EnableIpgroup *bool `json:"enable_ipgroup,omitempty"`

	// 访问控制组的类型。
	Type *string `json:"type,omitempty"`
}

func (o EsHealthIpgroupResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsHealthIpgroupResource struct{}"
	}

	return strings.Join([]string{"EsHealthIpgroupResource", string(data)}, " ")
}
