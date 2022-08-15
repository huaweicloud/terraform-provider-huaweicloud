package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateUserResult struct {

	// IAM用户访问方式。 - default：默认访问模式，编程访问和管理控制台访问。 - programmatic：编程访问。 - console：管理控制台访问。
	AccessMode *string `json:"access_mode,omitempty"`

	// IAM用户密码状态。true：需要修改密码，false：正常。
	PwdStatus *bool `json:"pwd_status,omitempty"`

	// IAM用户在外部系统中的ID。 >外部系统指与华为云对接的外部企业管理系统，xaccount_type、xaccount_id、xdomain_type、xdomain_id、xuser_type、xuser_id等参数值，无法在华为云获取，请咨询企业管理员。
	XuserId *string `json:"xuser_id,omitempty"`

	// IAM用户在外部系统中的类型。 >外部系统指与华为云对接的外部企业管理系统，xaccount_type、xaccount_id、xdomain_type、xdomain_id、xuser_type、xuser_id等参数值，无法在华为云获取，请咨询企业管理员。
	XuserType *string `json:"xuser_type,omitempty"`

	// IAM用户的新描述信息。
	Description *string `json:"description,omitempty"`

	// IAM用户新用户名，长度1~64之间，只能包含如下字符：大小写字母、空格、数字或特殊字符（-_.）且不能以数字开头。
	Name string `json:"name"`

	// IAM用户新手机号，纯数字，长度小于等于32字符。必须与国家码同时存在。
	Phone *string `json:"phone,omitempty"`

	// IAM用户所属账号ID。
	DomainId string `json:"domain_id"`

	// 是否启用IAM用户。true为启用，false为停用，默认为true。
	Enabled bool `json:"enabled"`

	// 国家码。中国大陆为“0086”。
	Areacode *string `json:"areacode,omitempty"`

	// IAM用户新邮箱。
	Email *string `json:"email,omitempty"`

	// IAM用户ID。
	Id string `json:"id"`

	Links *LinksSelf `json:"links"`

	// 密码过期时间（UTC时间），“null”表示密码不过期。
	PasswordExpiresAt *string `json:"password_expires_at,omitempty"`
}

func (o UpdateUserResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateUserResult struct{}"
	}

	return strings.Join([]string{"UpdateUserResult", string(data)}, " ")
}
