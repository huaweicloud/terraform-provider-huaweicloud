package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowTopUrlRequest Request Object
type ShowTopUrlRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 查询起始时间戳（单位：毫秒）。该时间戳的取值在转化为日期格式后须满足以下格式：XXXX-XX-XX 00:00:00
	StartTime int64 `json:"start_time"`

	// 查询结束时间戳（单位：毫秒）。该时间戳的取值在转化为日期格式后须满足以下格式：XXXX-XX-XX 00:00:00
	EndTime int64 `json:"end_time"`

	// 域名列表，多个域名以逗号（半角）分隔，如：www.test1.com,www.test2.com，all表示查询名下全部域名。如果域名在查询时间段内无数据，结果将不返回该域名的信息。
	DomainName string `json:"domain_name"`

	// mainland_china(中国大陆)，outside_mainland_china(中国大陆境外)，默认为global(全球)。
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
	GLOBAL                 ShowTopUrlRequestServiceArea
}

func GetShowTopUrlRequestServiceAreaEnum() ShowTopUrlRequestServiceAreaEnum {
	return ShowTopUrlRequestServiceAreaEnum{
		MAINLAND_CHINA: ShowTopUrlRequestServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: ShowTopUrlRequestServiceArea{
			value: "outside_mainland_china",
		},
		GLOBAL: ShowTopUrlRequestServiceArea{
			value: "global",
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
