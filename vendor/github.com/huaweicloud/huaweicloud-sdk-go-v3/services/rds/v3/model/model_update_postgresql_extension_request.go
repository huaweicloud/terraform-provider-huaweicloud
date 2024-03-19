package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlExtensionRequest Request Object
type UpdatePostgresqlExtensionRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ExtensionRequest `json:"body,omitempty"`
}

func (o UpdatePostgresqlExtensionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlExtensionRequest struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlExtensionRequest", string(data)}, " ")
}
