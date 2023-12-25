package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKafkaUserClientQuotaTaskResponse Response Object
type CreateKafkaUserClientQuotaTaskResponse struct {

	// 创建流控配置的任务ID
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateKafkaUserClientQuotaTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKafkaUserClientQuotaTaskResponse struct{}"
	}

	return strings.Join([]string{"CreateKafkaUserClientQuotaTaskResponse", string(data)}, " ")
}
