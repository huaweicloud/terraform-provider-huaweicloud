package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type IpAuthInfo struct {

	// 推流域名或播放域名
	Domain string `json:"domain"`

	// 鉴权类型。 包含如下取值： - WHITE：IP白名单鉴权。 - BLACK：IP黑名单鉴权。 - NONE：不鉴权。
	AuthType IpAuthInfoAuthType `json:"auth_type"`

	// IP黑名单列表，IP之间用;分隔，如192.168.0.0;192.168.0.8，最多支持配置100个IP。支持IP网段添加，例如127.0.0.1/24表示采用子网掩码中的前24位为有效位，即用32-24=8bit来表示主机号，该子网可以容纳2^8 - 2 = 254 台主机。故127.0.0.1/24 表示IP网段范围是：127.0.0.1~127.0.0.255
	IpAuthList string `json:"ip_auth_list"`
}

func (o IpAuthInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IpAuthInfo struct{}"
	}

	return strings.Join([]string{"IpAuthInfo", string(data)}, " ")
}

type IpAuthInfoAuthType struct {
	value string
}

type IpAuthInfoAuthTypeEnum struct {
	WHITE IpAuthInfoAuthType
	BLACK IpAuthInfoAuthType
	NONE  IpAuthInfoAuthType
}

func GetIpAuthInfoAuthTypeEnum() IpAuthInfoAuthTypeEnum {
	return IpAuthInfoAuthTypeEnum{
		WHITE: IpAuthInfoAuthType{
			value: "WHITE",
		},
		BLACK: IpAuthInfoAuthType{
			value: "BLACK",
		},
		NONE: IpAuthInfoAuthType{
			value: "NONE",
		},
	}
}

func (c IpAuthInfoAuthType) Value() string {
	return c.value
}

func (c IpAuthInfoAuthType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *IpAuthInfoAuthType) UnmarshalJSON(b []byte) error {
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
