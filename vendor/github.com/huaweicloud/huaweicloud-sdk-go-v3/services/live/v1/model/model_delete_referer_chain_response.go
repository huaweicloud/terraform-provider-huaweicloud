package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteRefererChainResponse Response Object
type DeleteRefererChainResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteRefererChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRefererChainResponse struct{}"
	}

	return strings.Join([]string{"DeleteRefererChainResponse", string(data)}, " ")
}
