package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SourceWithPort 域名信息。
type SourceWithPort struct {

	// 加速域名id。
	DomainId *string `json:"domain_id,omitempty"`

	// 源站IP（非内网IP）或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// 源站类型，ipaddr：源站IP、 domain：源站域名、obs_bucket：OBS桶域名。
	OriginType SourceWithPortOriginType `json:"origin_type"`

	// OBS桶类型。   - private: 私有桶（除桶ACL授权外的其他用户无桶的访问权限）。   - public: 公有桶（任何用户都可以对桶内对象进行读操作）。
	ObsBucketType *string `json:"obs_bucket_type,omitempty"`

	// 主备状态（1代表主源站；0代表备源站）。
	ActiveStandby int32 `json:"active_standby"`

	// 是否开OBS托管(0表示关闭,1表示则为开启)，源站类型为obs_bucket时传递。
	EnableObsWebHosting *int32 `json:"enable_obs_web_hosting,omitempty"`

	// HTTP端口，默认80
	HttpPort *int32 `json:"http_port,omitempty"`

	// HTTPS端口，默认443
	HttpsPort *int32 `json:"https_port,omitempty"`
}

func (o SourceWithPort) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceWithPort struct{}"
	}

	return strings.Join([]string{"SourceWithPort", string(data)}, " ")
}

type SourceWithPortOriginType struct {
	value string
}

type SourceWithPortOriginTypeEnum struct {
	IPADDR     SourceWithPortOriginType
	DOMAIN     SourceWithPortOriginType
	OBS_BUCKET SourceWithPortOriginType
}

func GetSourceWithPortOriginTypeEnum() SourceWithPortOriginTypeEnum {
	return SourceWithPortOriginTypeEnum{
		IPADDR: SourceWithPortOriginType{
			value: "ipaddr",
		},
		DOMAIN: SourceWithPortOriginType{
			value: "domain",
		},
		OBS_BUCKET: SourceWithPortOriginType{
			value: "obs_bucket",
		},
	}
}

func (c SourceWithPortOriginType) Value() string {
	return c.value
}

func (c SourceWithPortOriginType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceWithPortOriginType) UnmarshalJSON(b []byte) error {
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
