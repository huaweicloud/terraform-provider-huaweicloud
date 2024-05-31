package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SourcesRequestBody 源站信息
type SourcesRequestBody struct {

	// 加速域名id。
	DomainId *string `json:"domain_id,omitempty"`

	// 源站IP（非内网IP）或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// 源站类型取值：ipaddr：源站IP、 domain：源站域名、obs_bucket：OBS桶域名。
	OriginType SourcesRequestBodyOriginType `json:"origin_type"`

	// OBS桶类型，源站类型是“OBS桶域名”时需要传该参数，不传默认为“public”。   - private: 私有桶（除桶ACL授权外的其他用户无桶的访问权限）。   - public: 公有桶（任何用户都可以对桶内对象进行读操作）。
	ObsBucketType *string `json:"obs_bucket_type,omitempty"`

	// 主备状态，1代表主源站，0代表备源站。
	ActiveStandby int32 `json:"active_standby"`

	// 是否开启OBS静态网站托管(0表示关闭,1表示则为开启)，源站类型为obs_bucket时传递。
	EnableObsWebHosting *int32 `json:"enable_obs_web_hosting,omitempty"`
}

func (o SourcesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourcesRequestBody struct{}"
	}

	return strings.Join([]string{"SourcesRequestBody", string(data)}, " ")
}

type SourcesRequestBodyOriginType struct {
	value string
}

type SourcesRequestBodyOriginTypeEnum struct {
	IPADDR     SourcesRequestBodyOriginType
	DOMAIN     SourcesRequestBodyOriginType
	OBS_BUCKET SourcesRequestBodyOriginType
}

func GetSourcesRequestBodyOriginTypeEnum() SourcesRequestBodyOriginTypeEnum {
	return SourcesRequestBodyOriginTypeEnum{
		IPADDR: SourcesRequestBodyOriginType{
			value: "ipaddr",
		},
		DOMAIN: SourcesRequestBodyOriginType{
			value: "domain",
		},
		OBS_BUCKET: SourcesRequestBodyOriginType{
			value: "obs_bucket",
		},
	}
}

func (c SourcesRequestBodyOriginType) Value() string {
	return c.value
}

func (c SourcesRequestBodyOriginType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourcesRequestBodyOriginType) UnmarshalJSON(b []byte) error {
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
