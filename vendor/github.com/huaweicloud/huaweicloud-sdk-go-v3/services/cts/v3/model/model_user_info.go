package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UserInfo 用户信息。
type UserInfo struct {

	// 用户ID，参见《云审计服务API参考》“获取账号ID和项目ID”章节。
	Id *string `json:"id,omitempty"`

	// 用户名称。
	Name *string `json:"name,omitempty"`

	// 用户名称。
	UserName *string `json:"user_name,omitempty"`

	Domain *BaseUser `json:"domain,omitempty"`

	// 账号ID，参见《云审计服务API参考》“获取账号ID和项目ID”章节。
	AccountId *string `json:"account_id,omitempty"`

	// 访问密钥ID。
	AccessKeyId *string `json:"access_key_id,omitempty"`

	// 操作用户身份的 URN。 如果是 IAM 用户身份，格式如 iam::<account-id>:user:<user-name>。 如果是 IAM 委托会话 身份，格式如 sts::<account-id>:assumed-agency:<agency-name>/<agency-session-name>。 如果是 IAM 联邦身份，格式如 sts::<account-id>:external-user:<idp_id>/<user-session-name>。
	PrincipalUrn *string `json:"principal_urn,omitempty"`

	// 操作用户身份Id。 - 如果是 IAM 用户身份，格式为 <user-id>。 - 如果是 IAM 委托会话身份，格式为 <agency-id>:<agency-session-name>。 - 如果是 IAM 联邦身份，格式为 <idp_id>:<user-session-name>
	PrincipalId *string `json:"principal_id,omitempty"`

	// 是否是根用户。 - 值为“true”时，表示操作者是根用户。 - 值为“false”时，表示操作者是委托会话身份、联邦身份或非根用户的 IAM 用户。
	PrincipalIsRootUser *UserInfoPrincipalIsRootUser `json:"principal_is_root_user,omitempty"`

	// 操作者身份类型。
	Type *string `json:"type,omitempty"`

	// 发出请求的服务的名称。控制台操作时为[\"service.console\" ]
	InvokedBy *[]string `json:"invoked_by,omitempty"`

	SessionContext *SessionContext `json:"session_context,omitempty"`
}

func (o UserInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserInfo struct{}"
	}

	return strings.Join([]string{"UserInfo", string(data)}, " ")
}

type UserInfoPrincipalIsRootUser struct {
	value string
}

type UserInfoPrincipalIsRootUserEnum struct {
	TRUE  UserInfoPrincipalIsRootUser
	FALSE UserInfoPrincipalIsRootUser
}

func GetUserInfoPrincipalIsRootUserEnum() UserInfoPrincipalIsRootUserEnum {
	return UserInfoPrincipalIsRootUserEnum{
		TRUE: UserInfoPrincipalIsRootUser{
			value: "true",
		},
		FALSE: UserInfoPrincipalIsRootUser{
			value: "false",
		},
	}
}

func (c UserInfoPrincipalIsRootUser) Value() string {
	return c.value
}

func (c UserInfoPrincipalIsRootUser) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UserInfoPrincipalIsRootUser) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
