package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListInstancesResponse struct {

	// 实例列表
	Instances *[]ShowInstanceResp `json:"instances,omitempty"`

	// 实例数量。
	InstanceNum    *int32 `json:"instance_num,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListInstancesResponse", string(data)}, " ")
}
