package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateStorageModeRequest Request Object
type UpdateStorageModeRequest struct {
	Body *UpdateStorageModeReq `json:"body,omitempty"`
}

func (o UpdateStorageModeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateStorageModeRequest struct{}"
	}

	return strings.Join([]string{"UpdateStorageModeRequest", string(data)}, " ")
}
