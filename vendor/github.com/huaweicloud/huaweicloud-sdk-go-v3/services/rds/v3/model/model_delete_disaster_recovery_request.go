package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDisasterRecoveryRequest Request Object
type DeleteDisasterRecoveryRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *DeleteDisasterRecoveryRequestBody `json:"body,omitempty"`
}

func (o DeleteDisasterRecoveryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDisasterRecoveryRequest struct{}"
	}

	return strings.Join([]string{"DeleteDisasterRecoveryRequest", string(data)}, " ")
}
