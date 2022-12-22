package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AvailableResourceIdsInfo struct {

	// 资源ID
	ResourceId *string `json:"resource_id,omitempty"`

	// 当前时间
	CurrentTime *string `json:"current_time,omitempty"`

	// 是否共享配额   - shared：共享的   - unshared：非共享的
	SharedQuota *string `json:"shared_quota,omitempty"`
}

func (o AvailableResourceIdsInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AvailableResourceIdsInfo struct{}"
	}

	return strings.Join([]string{"AvailableResourceIdsInfo", string(data)}, " ")
}
