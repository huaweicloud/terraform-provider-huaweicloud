package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstanceTopicsRequest Request Object
type ListInstanceTopicsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ListInstanceTopicsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceTopicsRequest struct{}"
	}

	return strings.Join([]string{"ListInstanceTopicsRequest", string(data)}, " ")
}
