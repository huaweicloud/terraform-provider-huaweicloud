package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteRefererChainRequest Request Object
type DeleteRefererChainRequest struct {

	// 直播域名
	Domain string `json:"domain"`
}

func (o DeleteRefererChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteRefererChainRequest struct{}"
	}

	return strings.Join([]string{"DeleteRefererChainRequest", string(data)}, " ")
}
