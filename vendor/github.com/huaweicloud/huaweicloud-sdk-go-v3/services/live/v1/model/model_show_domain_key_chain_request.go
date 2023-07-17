package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainKeyChainRequest Request Object
type ShowDomainKeyChainRequest struct {

	// 直播域名，包括推流域名和播放域名
	Domain string `json:"domain"`
}

func (o ShowDomainKeyChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainKeyChainRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainKeyChainRequest", string(data)}, " ")
}
