package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateSinkTaskQuotaRequest struct {

	// 实例转储ID。 请参考[实例生命周期][查询实例]接口返回的数据。
	ConnectorId string `json:"connector_id"`

	Body *UpdateSinkTaskQuotaReq `json:"body,omitempty"`
}

func (o UpdateSinkTaskQuotaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSinkTaskQuotaRequest struct{}"
	}

	return strings.Join([]string{"UpdateSinkTaskQuotaRequest", string(data)}, " ")
}
