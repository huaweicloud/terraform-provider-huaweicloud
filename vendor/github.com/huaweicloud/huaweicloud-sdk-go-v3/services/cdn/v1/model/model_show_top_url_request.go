package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowTopUrlRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 查询起始时间戳（单位：毫秒）。该时间戳的取值在转化为日期格式后须满足以下格式：XXXX-XX-XX 00:00:00
	StartTime int64 `json:"start_time"`

	// 查询结束时间戳（单位：毫秒）。该时间戳的取值在转化为日期格式后须满足以下格式：XXXX-XX-XX 00:00:00
	EndTime int64 `json:"end_time"`

	// 域名列表，多个域名以逗号（半角）分隔，如：www.test1.com,www.test2.com ，ALL表示查询名下全部域名。
	DomainName string `json:"domain_name"`

	// mainland_china(中国大陆)，outside_mainland_china(中国大陆境外)，默认为mainland_china。
	ServiceArea *ShowTopUrlRequestServiceArea `json:"service_area,omitempty"`

	// 参数类型支持：flux(流量),req_num(请求总数)。
	StatType ShowTopUrlRequestStatType `json:"stat_type"`
}

func (o ShowTopUrlRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTopUrlRequest struct{}"
	}

	return strings.Join([]string{"ShowTopUrlRequest", string(data)}, " ")
}

type ShowTopUrlRequestServiceArea struct {
	value string
}

type ShowTopUrlRequestServiceAreaEnum struct {
	MAINLAND_CHINA         ShowTopUrlRequestServiceArea
	OUTSIDE_MAINLAND_CHINA ShowTopUrlRequestServiceArea
}

func GetShowTopUrlRequestServiceAreaEnum() ShowTopUrlRequestServiceAreaEnum {
	return ShowTopUrlRequestServiceAreaEnum{
		MAINLAND_CHINA: ShowTopUrlRequestServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: ShowTopUrlRequestServiceArea{
			value: "outside_mainland_china",
		},
	}
}

func (c ShowTopUrlRequestServiceArea) Value() string {
	return c.value
}

func (c ShowTopUrlRequestServiceArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTopUrlRequestServiceArea) UnmarshalJSON(b []byte) error {
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

type ShowTopUrlRequestStatType struct {
	value string
}

type ShowTopUrlRequestStatTypeEnum struct {
	FLUX    ShowTopUrlRequestStatType
	REQ_NUM ShowTopUrlRequestStatType
}

func GetShowTopUrlRequestStatTypeEnum() ShowTopUrlRequestStatTypeEnum {
	return ShowTopUrlRequestStatTypeEnum{
		FLUX: ShowTopUrlRequestStatType{
			value: "flux",
		},
		REQ_NUM: ShowTopUrlRequestStatType{
			value: "req_num",
		},
	}
}

func (c ShowTopUrlRequestStatType) Value() string {
	return c.value
}

func (c ShowTopUrlRequestStatType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowTopUrlRequestStatType) UnmarshalJSON(b []byte) error {
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
