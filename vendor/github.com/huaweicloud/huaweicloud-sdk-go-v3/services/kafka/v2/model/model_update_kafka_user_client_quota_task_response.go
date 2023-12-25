package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateKafkaUserClientQuotaTaskResponse Response Object
type UpdateKafkaUserClientQuotaTaskResponse struct {

	// 修改流控配置的任务ID
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateKafkaUserClientQuotaTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateKafkaUserClientQuotaTaskResponse struct{}"
	}

	return strings.Join([]string{"UpdateKafkaUserClientQuotaTaskResponse", string(data)}, " ")
}
