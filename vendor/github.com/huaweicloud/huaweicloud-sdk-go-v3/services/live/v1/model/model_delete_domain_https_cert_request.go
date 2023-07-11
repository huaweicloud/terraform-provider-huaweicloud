package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDomainHttpsCertRequest Request Object
type DeleteDomainHttpsCertRequest struct {

	// 直播播放域名
	Domain string `json:"domain"`
}

func (o DeleteDomainHttpsCertRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDomainHttpsCertRequest struct{}"
	}

	return strings.Join([]string{"DeleteDomainHttpsCertRequest", string(data)}, " ")
}
