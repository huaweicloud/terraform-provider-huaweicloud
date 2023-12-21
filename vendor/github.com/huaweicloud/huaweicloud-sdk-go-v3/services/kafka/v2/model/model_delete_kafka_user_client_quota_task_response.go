package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteKafkaUserClientQuotaTaskResponse Response Object
type DeleteKafkaUserClientQuotaTaskResponse struct {

	// 删除流控配置的任务ID
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteKafkaUserClientQuotaTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteKafkaUserClientQuotaTaskResponse struct{}"
	}

	return strings.Join([]string{"DeleteKafkaUserClientQuotaTaskResponse", string(data)}, " ")
}
