package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetRefererChainRequest Request Object
type SetRefererChainRequest struct {
	Body *SetRefererChainInfo `json:"body,omitempty"`
}

func (o SetRefererChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRefererChainRequest struct{}"
	}

	return strings.Join([]string{"SetRefererChainRequest", string(data)}, " ")
}
