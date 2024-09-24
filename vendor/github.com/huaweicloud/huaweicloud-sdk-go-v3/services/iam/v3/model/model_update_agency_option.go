package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateAgencyOption
type UpdateAgencyOption struct {

	// 被委托方账号ID。如果trust_domain_id和trust_domain_name都填写，则优先校验trust_domain_name。四个参数至少填写一个。
	TrustDomainId *string `json:"trust_domain_id,omitempty"`

	// 被委托方账号名。如果trust_domain_id和trust_domain_name都填写，则优先校验trust_domain_name。四个参数至少填写一个。
	TrustDomainName *string `json:"trust_domain_name,omitempty"`

	// 委托描述信息，长度不大于255位。四个参数至少填写一个。
	Description *string `json:"description,omitempty"`

	// 委托的期限，单位为“天”。默认为FOREVER。取值为“FOREVER\"表示委托的期限为永久，取值为\"ONEDAY\"表示委托的期限为一天,取值为自定义天数表示委托的期限为有限天数，如20。四个参数至少填写一个。
	Duration *string `json:"duration,omitempty"`
}

func (o UpdateAgencyOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAgencyOption struct{}"
	}

	return strings.Join([]string{"UpdateAgencyOption", string(data)}, " ")
}
