package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListLogLtsConfigsResponse Response Object
type ListLogLtsConfigsResponse struct {

	// 实例的LTS配置
	InstanceLtsConfigs *[]InstanceLtsConfigResp `json:"instance_lts_configs,omitempty"`

	// 结果集大小
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListLogLtsConfigsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogLtsConfigsResponse struct{}"
	}

	return strings.Join([]string{"ListLogLtsConfigsResponse", string(data)}, " ")
}
