package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type KeystoneShowSecurityComplianceByOptionRequest struct {

	// 待查询的账号ID，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	// 查询条件。该字段内容为：password_regex或password_regex_description。  password_regex：密码强度策略的正则表达式；password_regex_description：密码强度策略的描述。
	Option KeystoneShowSecurityComplianceByOptionRequestOption `json:"option"`
}

func (o KeystoneShowSecurityComplianceByOptionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowSecurityComplianceByOptionRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowSecurityComplianceByOptionRequest", string(data)}, " ")
}

type KeystoneShowSecurityComplianceByOptionRequestOption struct {
	value string
}

type KeystoneShowSecurityComplianceByOptionRequestOptionEnum struct {
	PASSWORD_REGEX             KeystoneShowSecurityComplianceByOptionRequestOption
	PASSWORD_REGEX_DESCRIPTION KeystoneShowSecurityComplianceByOptionRequestOption
}

func GetKeystoneShowSecurityComplianceByOptionRequestOptionEnum() KeystoneShowSecurityComplianceByOptionRequestOptionEnum {
	return KeystoneShowSecurityComplianceByOptionRequestOptionEnum{
		PASSWORD_REGEX: KeystoneShowSecurityComplianceByOptionRequestOption{
			value: "password_regex",
		},
		PASSWORD_REGEX_DESCRIPTION: KeystoneShowSecurityComplianceByOptionRequestOption{
			value: "password_regex_description",
		},
	}
}

func (c KeystoneShowSecurityComplianceByOptionRequestOption) Value() string {
	return c.value
}

func (c KeystoneShowSecurityComplianceByOptionRequestOption) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *KeystoneShowSecurityComplianceByOptionRequestOption) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
