package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDomainHttpsCertRequest Request Object
type UpdateDomainHttpsCertRequest struct {

	// 直播播放域名
	Domain string `json:"domain"`

	Body *DomainHttpsCertInfo `json:"body,omitempty"`
}

func (o UpdateDomainHttpsCertRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainHttpsCertRequest struct{}"
	}

	return strings.Join([]string{"UpdateDomainHttpsCertRequest", string(data)}, " ")
}
