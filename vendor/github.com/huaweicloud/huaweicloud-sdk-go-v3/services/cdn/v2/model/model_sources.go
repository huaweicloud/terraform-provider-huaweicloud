package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Sources 源站信息。
type Sources struct {

	// 加速域名id。
	DomainId *string `json:"domain_id,omitempty"`

	// 源站IP（非内网IP）或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// 源站类型取值：ipaddr：源站IP、 domain：源站域名、obs_bucket：OBS桶域名。
	OriginType SourcesOriginType `json:"origin_type"`

	// OBS桶类型。   - private: 私有桶（除桶ACL授权外的其他用户无桶的访问权限）。   - public: 公有桶（任何用户都可以对桶内对象进行读操作）。
	ObsBucketType *string `json:"obs_bucket_type,omitempty"`

	// 主备状态，1代表主源站，0代表备源站。
	ActiveStandby int32 `json:"active_standby"`

	// 是否开启OBS静态网站托管(0表示关闭,1表示则为开启)，源站类型为obs_bucket时传递。
	EnableObsWebHosting *int32 `json:"enable_obs_web_hosting,omitempty"`
}

func (o Sources) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Sources struct{}"
	}

	return strings.Join([]string{"Sources", string(data)}, " ")
}

type SourcesOriginType struct {
	value string
}

type SourcesOriginTypeEnum struct {
	IPADDR     SourcesOriginType
	DOMAIN     SourcesOriginType
	OBS_BUCKET SourcesOriginType
}

func GetSourcesOriginTypeEnum() SourcesOriginTypeEnum {
	return SourcesOriginTypeEnum{
		IPADDR: SourcesOriginType{
			value: "ipaddr",
		},
		DOMAIN: SourcesOriginType{
			value: "domain",
		},
		OBS_BUCKET: SourcesOriginType{
			value: "obs_bucket",
		},
	}
}

func (c SourcesOriginType) Value() string {
	return c.value
}

func (c SourcesOriginType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourcesOriginType) UnmarshalJSON(b []byte) error {
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
