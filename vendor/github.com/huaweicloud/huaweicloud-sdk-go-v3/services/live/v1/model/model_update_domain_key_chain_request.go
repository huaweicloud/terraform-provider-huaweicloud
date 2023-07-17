package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainKeyChainRequest Request Object
type UpdateDomainKeyChainRequest struct {

	// 直播域名，包括推流域名和播放域名
	Domain string `json:"domain"`

	Body *KeyChainInfo `json:"body,omitempty"`
}

func (o UpdateDomainKeyChainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainKeyChainRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainKeyChainRequest", string(data)}, " ")
}
