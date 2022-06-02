package model

import (
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
	"strings"
)

// Response Object
type CreateDomainResponse struct {

	// 直播域名
	Domain *string `json:"domain,omitempty"`

	// 域名类型 - pull表示播放域名 - push表示推流域名
	DomainType *CreateDomainResponseDomainType `json:"domain_type,omitempty"`

	// 直播域名的CNAME
	DomainCname *string `json:"domain_cname,omitempty"`

	// 直播所属直播中心
	Region *string `json:"region,omitempty"`

	// 直播域名的状态
	Status *CreateDomainResponseStatus `json:"status,omitempty"`

	// 域名创建时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间
	CreateTime *sdktime.SdkTime `json:"create_time,omitempty"`

	// 状态描述
	StatusDescribe *string `json:"status_describe,omitempty"`

	// 域名应用区域 - mainland_china表示中国大陆区域 - outside_mainland_china表示中国大陆以外区域 - global表示全球区域
	ServiceArea    *CreateDomainResponseServiceArea `json:"service_area,omitempty"`
	HttpStatusCode int                              `json:"-"`
}

func (o CreateDomainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDomainResponse struct{}"
	}

	return strings.Join([]string{"CreateDomainResponse", string(data)}, " ")
}

type CreateDomainResponseDomainType struct {
	value string
}

type CreateDomainResponseDomainTypeEnum struct {
	PULL CreateDomainResponseDomainType
	PUSH CreateDomainResponseDomainType
}

func GetCreateDomainResponseDomainTypeEnum() CreateDomainResponseDomainTypeEnum {
	return CreateDomainResponseDomainTypeEnum{
		PULL: CreateDomainResponseDomainType{
			value: "pull",
		},
		PUSH: CreateDomainResponseDomainType{
			value: "push",
		},
	}
}

func (c CreateDomainResponseDomainType) Value() string {
	return c.value
}

func (c CreateDomainResponseDomainType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateDomainResponseDomainType) UnmarshalJSON(b []byte) error {
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

type CreateDomainResponseStatus struct {
	value string
}

type CreateDomainResponseStatusEnum struct {
	ON          CreateDomainResponseStatus
	OFF         CreateDomainResponseStatus
	CONFIGURING CreateDomainResponseStatus
	DISABLE     CreateDomainResponseStatus
}

func GetCreateDomainResponseStatusEnum() CreateDomainResponseStatusEnum {
	return CreateDomainResponseStatusEnum{
		ON: CreateDomainResponseStatus{
			value: "on",
		},
		OFF: CreateDomainResponseStatus{
			value: "off",
		},
		CONFIGURING: CreateDomainResponseStatus{
			value: "configuring",
		},
		DISABLE: CreateDomainResponseStatus{
			value: "disable",
		},
	}
}

func (c CreateDomainResponseStatus) Value() string {
	return c.value
}

func (c CreateDomainResponseStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateDomainResponseStatus) UnmarshalJSON(b []byte) error {
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

type CreateDomainResponseServiceArea struct {
	value string
}

type CreateDomainResponseServiceAreaEnum struct {
	MAINLAND_CHINA         CreateDomainResponseServiceArea
	OUTSIDE_MAINLAND_CHINA CreateDomainResponseServiceArea
	GLOBAL                 CreateDomainResponseServiceArea
}

func GetCreateDomainResponseServiceAreaEnum() CreateDomainResponseServiceAreaEnum {
	return CreateDomainResponseServiceAreaEnum{
		MAINLAND_CHINA: CreateDomainResponseServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: CreateDomainResponseServiceArea{
			value: "outside_mainland_china",
		},
		GLOBAL: CreateDomainResponseServiceArea{
			value: "global",
		},
	}
}

func (c CreateDomainResponseServiceArea) Value() string {
	return c.value
}

func (c CreateDomainResponseServiceArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateDomainResponseServiceArea) UnmarshalJSON(b []byte) error {
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
