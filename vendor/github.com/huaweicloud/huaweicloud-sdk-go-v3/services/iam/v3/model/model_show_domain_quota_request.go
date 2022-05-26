package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowDomainQuotaRequest struct {

	// 待查询的账号ID，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId string `json:"domain_id"`

	// 查询配额的类型，取值范围为：user, group, idp, agency, policy, assigment_group_mp, assigment_agency_mp, assigment_group_ep, assigment_user_ep。
	Type *ShowDomainQuotaRequestType `json:"type,omitempty"`
}

func (o ShowDomainQuotaRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainQuotaRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainQuotaRequest", string(data)}, " ")
}

type ShowDomainQuotaRequestType struct {
	value string
}

type ShowDomainQuotaRequestTypeEnum struct {
	USER                ShowDomainQuotaRequestType
	GROUP               ShowDomainQuotaRequestType
	IDP                 ShowDomainQuotaRequestType
	AGENCY              ShowDomainQuotaRequestType
	POLICY              ShowDomainQuotaRequestType
	ASSIGMENT_GROUP_MP  ShowDomainQuotaRequestType
	ASSIGMENT_AGENCY_MP ShowDomainQuotaRequestType
	ASSIGMENT_GROUP_EP  ShowDomainQuotaRequestType
	ASSIGMENT_USER_EP   ShowDomainQuotaRequestType
}

func GetShowDomainQuotaRequestTypeEnum() ShowDomainQuotaRequestTypeEnum {
	return ShowDomainQuotaRequestTypeEnum{
		USER: ShowDomainQuotaRequestType{
			value: "user",
		},
		GROUP: ShowDomainQuotaRequestType{
			value: "group",
		},
		IDP: ShowDomainQuotaRequestType{
			value: "idp",
		},
		AGENCY: ShowDomainQuotaRequestType{
			value: "agency",
		},
		POLICY: ShowDomainQuotaRequestType{
			value: "policy",
		},
		ASSIGMENT_GROUP_MP: ShowDomainQuotaRequestType{
			value: "assigment_group_mp",
		},
		ASSIGMENT_AGENCY_MP: ShowDomainQuotaRequestType{
			value: "assigment_agency_mp",
		},
		ASSIGMENT_GROUP_EP: ShowDomainQuotaRequestType{
			value: "assigment_group_ep",
		},
		ASSIGMENT_USER_EP: ShowDomainQuotaRequestType{
			value: "assigment_user_ep",
		},
	}
}

func (c ShowDomainQuotaRequestType) Value() string {
	return c.value
}

func (c ShowDomainQuotaRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowDomainQuotaRequestType) UnmarshalJSON(b []byte) error {
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
