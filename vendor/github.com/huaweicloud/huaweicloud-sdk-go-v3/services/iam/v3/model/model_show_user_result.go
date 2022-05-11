package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ShowUserResult struct {

	// IAM用户是否启用。true表示启用，false表示停用，默认为true。
	Enabled bool `json:"enabled"`

	// IAM用户ID。
	Id string `json:"id"`

	// IAM用户所属账号ID。
	DomainId string `json:"domain_id"`

	// IAM用户名。
	Name string `json:"name"`

	Links *Links `json:"links"`

	// IAM用户在外部系统中的ID。
	XuserId *string `json:"xuser_id,omitempty"`

	// IAM用户在外部系统中的类型。
	XuserType *string `json:"xuser_type,omitempty"`

	// IAM用户手机号的国家码。
	Areacode *string `json:"areacode,omitempty"`

	// IAM用户邮箱。
	Email *string `json:"email,omitempty"`

	// IAM用户手机号。
	Phone *string `json:"phone,omitempty"`

	// IAM用户密码状态。true：需要修改密码，false：正常。
	PwdStatus *bool `json:"pwd_status,omitempty"`

	// IAM用户更新时间。
	UpdateTime *string `json:"update_time,omitempty"`

	// IAM用户创建时间。
	CreateTime *string `json:"create_time,omitempty"`

	// IAM用户最后登录时间。
	LastLoginTime *string `json:"last_login_time,omitempty"`

	// IAM用户密码强度。结果为Low/Middle/High/None，分别表示密码强度低/中/高/无。
	PwdStrength *string `json:"pwd_strength,omitempty"`

	// IAM用户是否为根用户。
	IsDomainOwner bool `json:"is_domain_owner"`

	// IAM用户访问模式。
	AccessMode string `json:"access_mode"`

	// IAM用户描述信息
	Description string `json:"description"`
}

func (o ShowUserResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUserResult struct{}"
	}

	return strings.Join([]string{"ShowUserResult", string(data)}, " ")
}
