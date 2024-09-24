package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetRefererChainResponse Response Object
type SetRefererChainResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetRefererChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRefererChainResponse struct{}"
	}

	return strings.Join([]string{"SetRefererChainResponse", string(data)}, " ")
}
