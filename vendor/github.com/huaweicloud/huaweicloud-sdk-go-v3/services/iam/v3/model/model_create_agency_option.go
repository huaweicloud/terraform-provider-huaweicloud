package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAgencyOption
type CreateAgencyOption struct {

	// 委托名，长度不大于64位。
	Name string `json:"name"`

	// 委托方账号ID。
	DomainId string `json:"domain_id"`

	// 被委托方账号ID。trust_domain_id和trust_domain_name至少填写一个，若都填写优先校验trust_domain_name。
	TrustDomainId *string `json:"trust_domain_id,omitempty"`

	// 被委托方账号名。trust_domain_id和trust_domain_name至少填写一个，若都填写优先校验trust_domain_name。
	TrustDomainName *string `json:"trust_domain_name,omitempty"`

	// 委托描述信息，长度不大于255位。
	Description *string `json:"description,omitempty"`

	// description: 委托的期限，单位为“天”。默认为FOREVER。取值为“FOREVER\"表示委托的期限为永久，取值为\"ONEDAY\"表示委托的期限为一天,取值为自定义天数表示委托的期限为有限天数，如20。四个参数至少填写一个。
	Duration *string `json:"duration,omitempty"`
}

func (o CreateAgencyOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAgencyOption struct{}"
	}

	return strings.Join([]string{"CreateAgencyOption", string(data)}, " ")
}
