package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchShutdownInsReq struct {

	// 实例id列表
	InstanceIds []string `json:"instance_ids"`
}

func (o BatchShutdownInsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchShutdownInsReq struct{}"
	}

	return strings.Join([]string{"BatchShutdownInsReq", string(data)}, " ")
}
