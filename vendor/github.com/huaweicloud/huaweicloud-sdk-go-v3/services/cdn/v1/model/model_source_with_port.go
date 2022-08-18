package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源站信息
type SourceWithPort struct {

	// 源站IP（非内网IP）或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// 源站类型（\"ipaddr\"： \"IP源站\"；\"domain\"： \"域名源站\"；\"obs_bucket\"： \"OBS Bucket源站\"）
	OriginType SourceWithPortOriginType `json:"origin_type"`

	// 主备状态（1代表主站；0代表备站）；主源站必须存在，备源站可选。
	ActiveStandby int32 `json:"active_standby"`

	// 是否开启Obs静态网站托管(0表示关闭,1表示则为开启)，源站类型为obs_bucket时传递。
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
