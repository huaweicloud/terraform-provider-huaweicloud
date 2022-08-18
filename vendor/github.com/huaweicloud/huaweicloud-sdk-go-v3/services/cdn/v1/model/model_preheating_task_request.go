package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PreheatingTaskRequest struct {
	PreheatingTask *PreheatingTaskRequestBody `json:"preheating_task"`
}

func (o PreheatingTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PreheatingTaskRequest struct{}"
	}

	return strings.Join([]string{"PreheatingTaskRequest", string(data)}, " ")
}
