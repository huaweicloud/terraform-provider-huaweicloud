package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVerifyDomainOwnerInfoRequest Request Object
type ShowVerifyDomainOwnerInfoRequest struct {

	// 域名
	DomainName string `json:"domain_name"`
}

func (o ShowVerifyDomainOwnerInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVerifyDomainOwnerInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowVerifyDomainOwnerInfoRequest", string(data)}, " ")
}
