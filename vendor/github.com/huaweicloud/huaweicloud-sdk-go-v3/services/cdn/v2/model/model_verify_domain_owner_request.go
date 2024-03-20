package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VerifyDomainOwnerRequest Request Object
type VerifyDomainOwnerRequest struct {

	// 域名
	DomainName string `json:"domain_name"`

	Body *VerifyDomainOwnerRequestBody `json:"body,omitempty"`
}

func (o VerifyDomainOwnerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VerifyDomainOwnerRequest struct{}"
	}

	return strings.Join([]string{"VerifyDomainOwnerRequest", string(data)}, " ")
}
