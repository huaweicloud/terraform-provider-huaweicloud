package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainHttpsCertRequest Request Object
type ShowDomainHttpsCertRequest struct {

	// 直播播放域名
	Domain string `json:"domain"`
}

func (o ShowDomainHttpsCertRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainHttpsCertRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainHttpsCertRequest", string(data)}, " ")
}
