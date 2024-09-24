package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRefererChainRequest Request Object
type ShowRefererChainRequest struct {

	// 直播播放域名
	Domain string `json:"domain"`
}

func (o ShowRefererChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRefererChainRequest struct{}"
	}

	return strings.Join([]string{"ShowRefererChainRequest", string(data)}, " ")
}
