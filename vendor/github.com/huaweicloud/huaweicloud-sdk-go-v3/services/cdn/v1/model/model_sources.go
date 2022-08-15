package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源站信息
type Sources struct {

	// 源站IP（非内网IP）或者域名。
	IpOrDomain string `json:"ip_or_domain"`

	// 源站类型取值：ipaddr、 domain、obs_bucket，分别表示：源站IP、源站域名、OBS桶访问域名。
	OriginType SourcesOriginType `json:"origin_type"`

	// 主备状态（1代表主站；0代表备站）,主源站必须存在，备源站可选，OBS桶不能有备源站。
	ActiveStandby int32 `json:"active_standby"`

	// 是否开启Obs静态网站托管(0表示关闭,1表示则为开启)，源站类型为obs_bucket时传递。
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
