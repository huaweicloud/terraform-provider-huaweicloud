package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateHlsConfigRequest Request Object
type UpdateHlsConfigRequest struct {
	Body *ModifyHlsConfig `json:"body,omitempty"`
}

func (o UpdateHlsConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateHlsConfigRequest struct{}"
	}

	return strings.Join([]string{"UpdateHlsConfigRequest", string(data)}, " ")
}
