package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowPermanentAccessKeyRequest struct {

	// 待查询的指定AK。
	AccessKey string `json:"access_key"`
}

func (o ShowPermanentAccessKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPermanentAccessKeyRequest struct{}"
	}

	return strings.Join([]string{"ShowPermanentAccessKeyRequest", string(data)}, " ")
}
