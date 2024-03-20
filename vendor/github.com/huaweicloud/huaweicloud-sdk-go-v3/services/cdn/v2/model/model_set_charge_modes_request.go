package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetChargeModesRequest Request Object
type SetChargeModesRequest struct {
	Body *SetChargeModesBody `json:"body,omitempty"`
}

func (o SetChargeModesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetChargeModesRequest struct{}"
	}

	return strings.Join([]string{"SetChargeModesRequest", string(data)}, " ")
}
