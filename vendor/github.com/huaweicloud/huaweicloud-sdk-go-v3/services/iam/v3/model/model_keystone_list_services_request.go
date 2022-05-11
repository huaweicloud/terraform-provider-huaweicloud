package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListServicesRequest struct {

	// 服务类型。
	Type *string `json:"type,omitempty"`
}

func (o KeystoneListServicesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListServicesRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListServicesRequest", string(data)}, " ")
}
