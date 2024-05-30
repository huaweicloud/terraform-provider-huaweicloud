package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DomainBody 创建域名Body
type DomainBody struct {

	// 加速域名。（ 由字母（A-Z，a-z，大小写等价）、数字（0-9）和连接符（-）组成，各级域名之间用（.）连接，域名长度不超过75个字符。连接符（-）不能作为域名的开头或结尾字符。）
	DomainName string `json:"domain_name"`

	// 域名业务类型，若为web，则表示类型为网页加速；若为download，则表示业务类型为文件下载加速；若为video，则表示业务类型为点播加速；若为wholeSite，则表示业务类型为全站加速。
	BusinessType DomainBodyBusinessType `json:"business_type"`

	// 源站配置。
	Sources []SourcesRequestBody `json:"sources"`

	// 域名服务范围，若为mainland_china，则表示服务范围为中国大陆；若为outside_mainland_china，则表示服务范围为中国大陆境外；若为global，则表示服务范围为全球。
	ServiceArea DomainBodyServiceArea `json:"service_area"`

	// 当用户开启企业项目功能时，该参数生效，表示添加加速域名到该企业项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o DomainBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainBody struct{}"
	}

	return strings.Join([]string{"DomainBody", string(data)}, " ")
}

type DomainBodyBusinessType struct {
	value string
}

type DomainBodyBusinessTypeEnum struct {
	WEB        DomainBodyBusinessType
	DOWNLOAD   DomainBodyBusinessType
	VIDEO      DomainBodyBusinessType
	WHOLE_SITE DomainBodyBusinessType
}

func GetDomainBodyBusinessTypeEnum() DomainBodyBusinessTypeEnum {
	return DomainBodyBusinessTypeEnum{
		WEB: DomainBodyBusinessType{
			value: "web",
		},
		DOWNLOAD: DomainBodyBusinessType{
			value: "download",
		},
		VIDEO: DomainBodyBusinessType{
			value: "video",
		},
		WHOLE_SITE: DomainBodyBusinessType{
			value: "wholeSite",
		},
	}
}

func (c DomainBodyBusinessType) Value() string {
	return c.value
}

func (c DomainBodyBusinessType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DomainBodyBusinessType) UnmarshalJSON(b []byte) error {
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

type DomainBodyServiceArea struct {
	value string
}

type DomainBodyServiceAreaEnum struct {
	MAINLAND_CHINA         DomainBodyServiceArea
	OUTSIDE_MAINLAND_CHINA DomainBodyServiceArea
	GLOBAL                 DomainBodyServiceArea
}

func GetDomainBodyServiceAreaEnum() DomainBodyServiceAreaEnum {
	return DomainBodyServiceAreaEnum{
		MAINLAND_CHINA: DomainBodyServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: DomainBodyServiceArea{
			value: "outside_mainland_china",
		},
		GLOBAL: DomainBodyServiceArea{
			value: "global",
		},
	}
}

func (c DomainBodyServiceArea) Value() string {
	return c.value
}

func (c DomainBodyServiceArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DomainBodyServiceArea) UnmarshalJSON(b []byte) error {
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
