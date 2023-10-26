package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSinkTaskRequest Request Object
type CreateSinkTaskRequest struct {

	// 实例转储ID。  请参考[查询实例](ShowInstance.xml)返回的数据。
	ConnectorId string `json:"connector_id"`

	Body *CreateSinkTaskReq `json:"body,omitempty"`
}

func (o CreateSinkTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSinkTaskRequest struct{}"
	}

	return strings.Join([]string{"CreateSinkTaskRequest", string(data)}, " ")
}
