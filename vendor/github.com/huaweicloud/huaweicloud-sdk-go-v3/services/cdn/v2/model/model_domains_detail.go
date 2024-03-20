package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DomainsDetail 域名信息。
type DomainsDetail struct {

	// 加速域名ID。
	Id *string `json:"id,omitempty"`

	// 加速域名。
	DomainName *string `json:"domain_name,omitempty"`

	// 域名业务类型，若为web，则表示类型为网站加速；若为download，则表示业务类型为文件下载加速；若为video，则表示业务类型为点播加速；若为wholeSite，则表示类型为全站加速。
	BusinessType *string `json:"business_type,omitempty"`

	// 加速域名状态。取值意义： - online表示“已开启” - offline表示“已停用” - configuring表示“配置中” - configure_failed表示“配置失败” - checking表示“审核中” - check_failed表示“审核未通过” - deleting表示“删除中”。
	DomainStatus *string `json:"domain_status,omitempty"`

	// 加速域名对应的CNAME。
	Cname *string `json:"cname,omitempty"`

	// 源站配置。
	Sources *[]SourcesDomainConfig `json:"sources,omitempty"`

	// 是否开启HTTPS加速。 0：代表未开启HTTPS加速； 1：代表开启HTTPS加速，且回源方式为协议跟随； 2：代表开启HTTPS加速，且回源方式为HTTP； 3：代表开启HTTPS加速，且回源方式为HTTPS。
	HttpsStatus *int32 `json:"https_status,omitempty"`

	// 域名创建时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 域名修改时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	UpdateTime *int64 `json:"update_time,omitempty"`

	// 封禁状态（0代表未禁用；1代表禁用）。
	Disabled *int32 `json:"disabled,omitempty"`

	// 锁定状态（0代表未锁定；1代表锁定）。
	Locked *int32 `json:"locked,omitempty"`

	// 华为云CDN提供的加速服务范围，包含：mainland_china中国大陆、outside_mainland_china中国大陆境外、global全球。
	ServiceArea *DomainsDetailServiceArea `json:"service_area,omitempty"`
}

func (o DomainsDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainsDetail struct{}"
	}

	return strings.Join([]string{"DomainsDetail", string(data)}, " ")
}

type DomainsDetailServiceArea struct {
	value string
}

type DomainsDetailServiceAreaEnum struct {
	MAINLAND_CHINA         DomainsDetailServiceArea
	OUTSIDE_MAINLAND_CHINA DomainsDetailServiceArea
	GLOBAL                 DomainsDetailServiceArea
}

func GetDomainsDetailServiceAreaEnum() DomainsDetailServiceAreaEnum {
	return DomainsDetailServiceAreaEnum{
		MAINLAND_CHINA: DomainsDetailServiceArea{
			value: "mainland_china",
		},
		OUTSIDE_MAINLAND_CHINA: DomainsDetailServiceArea{
			value: "outside_mainland_china",
		},
		GLOBAL: DomainsDetailServiceArea{
			value: "global",
		},
	}
}

func (c DomainsDetailServiceArea) Value() string {
	return c.value
}

func (c DomainsDetailServiceArea) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DomainsDetailServiceArea) UnmarshalJSON(b []byte) error {
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
