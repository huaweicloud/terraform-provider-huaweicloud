package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdatePermanentAccessKeyRequest struct {

	// 待修改的指定AK。
	AccessKey string `json:"access_key"`

	Body *UpdatePermanentAccessKeyRequestBody `json:"body,omitempty"`
}

func (o UpdatePermanentAccessKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePermanentAccessKeyRequest struct{}"
	}

	return strings.Join([]string{"UpdatePermanentAccessKeyRequest", string(data)}, " ")
}
