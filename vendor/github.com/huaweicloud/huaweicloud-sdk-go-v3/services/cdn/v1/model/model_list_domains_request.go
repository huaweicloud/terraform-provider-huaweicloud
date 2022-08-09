package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListDomainsRequest struct {

	// 加速域名，采用模糊匹配的方式。（长度限制为1-255字符）。
	DomainName *string `json:"domain_name,omitempty"`

	// 加速域名的业务类型。取值： - web（网站加速） - download（文件下载加速） - video（点播加速） - wholeSite（全站加速）
	BusinessType *ListDomainsRequestBusinessType `json:"business_type,omitempty"`

	// 加速域名状态。取值意义： - online表示“已开启” - offline表示“已停用” - configuring表示“配置中” - configure_failed表示“配置失败” - checking表示“审核中” - check_failed表示“审核未通过” - deleting表示“删除中”。
	DomainStatus *ListDomainsRequestDomainStatus `json:"domain_status,omitempty"`

	// 华为云CDN提供的加速服务范围，包含： - mainland_china 中国大陆 - outside_mainland_china 中国大陆境外 - global 全球。
	ServiceArea *ListDomainsRequestServiceArea `json:"service_area,omitempty"`

	// 每页的数量，取值范围1-10000，不设值时默认值为30。
	PageSize *int32 `json:"page_size,omitempty"`

	// 查询的页码。取值范围1-65535，不设值时默认值为1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
}

func (o ListDomainsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDomainsRequest struct{}"
	}

	return strings.Join([]string{"ListDomainsRequest", string(data)}, " ")
}

type ListDomainsRequestBusinessType struct {
	value string
}

type ListDomainsRequestBusinessTypeEnum struct {
	WEB        ListDomainsRequestBusinessType
	DOWNLOAD   ListDomainsRequestBusinessType
	VIDEO      ListDomainsRequestBusinessType
	WHOLE_SITE ListDomainsRequestBusinessType
}

func GetListDomainsRequestBusinessTypeEnum() ListDomainsRequestBusinessTypeEnum {
	return ListDomainsRequestBusinessTypeEnum{
		WEB: ListDomainsRequestBusinessType{
			value: "web",
		},
		DOWNLOAD: ListDomainsRequestBusinessType{
			value: "download",
		},
		VIDEO: ListDomainsRequestBusinessType{
			value: "video",
		},
		WHOLE_SITE: ListDomainsRequestBusinessType{
			value: "wholeSite",
		},
	}
}

func (c ListDomainsRequestBusinessType) Value() string {
	return c.value
}

func (c ListDomainsRequestBusinessType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListDomainsRequestBusinessType) UnmarshalJSON(b []byte) error {
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

type ListDomainsRequestDomainStatus struct {
	value string
}

type ListDomainsRequestDomainStatusEnum struct {
	ONLINE           ListDomainsRequestDomainStatus
	OFFLINE          ListDomainsRequestDomainStatus
	CONFIGURING      ListDomainsRequestDomainStatus
	CONFIGURE_FAILED ListDomainsRequestDomainStatus
	CHECKING         ListDomainsRequestDomainStatus
	CHECK_FAILED     ListDomainsRequestDomainStatus
	DELETING         ListDomainsRequestDomainStatus
}

func GetListDomainsRequestDomainStatusEnum() ListDomainsRequestDomainStatusEnum {
	return ListDomainsRequestDomainStatusEnum{
		ONLINE: ListDomainsRequestDomainStatus{
			value: "online",
		},
		OFFLINE: ListDomainsRequestDomainStatus{
			value: "offline",
		},
		CONFIGURING: ListDomainsRequestDomainStatus{
			value: "configuring",
		},
		CONFIGURE_FAILED: ListDomainsRequestDomainStatus{
			value: "configure_failed",
		},
		CHECKING: ListDomainsRequestDomainStatus{
			value: "checking",
		},
		CHECK_FAILED: ListDomainsRequestDomainStatus{
			value: "check_failed",
		},
		DELETING: ListDomainsRequestDomainStatus{
			value: "deleting",
		},
	}
}

func (c ListDomainsRequestDomainStatus) Value() string {
	return c.value
}

func (c ListDomainsRequestDomainStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListDomainsRequestDomainStatus) UnmarshalJSON(b []byte) error {
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

type ListDomainsRequestServiceArea struct {
	value string
}

type ListDomainsRequestServiceAreaEnum struct {
	MAINLAND_CHINA         ListDomainsRequestServiceArea
	OUTSIDE_MAINLAND_CHINA ListDomainsRequestServiceArea
	GLOBAL                 ListDomainsRequestServiceArea
}

func GetListDomainsRequestServiceAreaEnum() ListDomainsRequestServiceAreaEnum {
	return ListDomainsRequestServiceAreaEnum{
		MAINLAND_CHINA: ListDomainsRequestServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: ListDomainsRequestServiceArea{
			value: "outside_mainland_china",
		},
		GLOBAL: ListDomainsRequestServiceArea{
			value: "global",
		},
	}
}

func (c ListDomainsRequestServiceArea) Value() string {
	return c.value
}

func (c ListDomainsRequestServiceArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListDomainsRequestServiceArea) UnmarshalJSON(b []byte) error {
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
