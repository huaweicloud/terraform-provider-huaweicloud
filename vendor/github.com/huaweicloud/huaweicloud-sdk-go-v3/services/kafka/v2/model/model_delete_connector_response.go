package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConnectorResponse Response Object
type DeleteConnectorResponse struct {

	// 返回异步执行删除任务的job id。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteConnectorResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConnectorResponse struct{}"
	}

	return strings.Join([]string{"DeleteConnectorResponse", string(data)}, " ")
}
