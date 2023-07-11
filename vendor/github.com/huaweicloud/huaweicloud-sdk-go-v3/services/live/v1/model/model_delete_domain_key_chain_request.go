package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDomainKeyChainRequest Request Object
type DeleteDomainKeyChainRequest struct {

	// 直播域名，包括推流域名和播放域名
	Domain string `json:"domain"`
}

func (o DeleteDomainKeyChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainKeyChainRequest struct{}"
	}

	return strings.Join([]string{"DeleteDomainKeyChainRequest", string(data)}, " ")
}
