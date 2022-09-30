package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateConnectorResponse struct {

	// 任务ID。
	JobId *string `json:"job_id,omitempty"`

	// 实例转储ID。
	ConnectorId    *string `json:"connector_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateConnectorResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateConnectorResponse struct{}"
	}

	return strings.Join([]string{"CreateConnectorResponse", string(data)}, " ")
}
