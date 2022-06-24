package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源端CDN配置返回值。
type SourceCdnResp struct {

	//   从指定域名获取对象。
	Domain string `json:"domain"`

	// 协议类型，支持http和https协议。
	Protocol SourceCdnRespProtocol `json:"protocol"`

	// 鉴权类型: NONE, QINIU_PRIVATE_AUTHENTICATION, ALIYUN_OSS_A, ALIYUN_OSS_B, ALIYUN_OSS_C, KSYUN_PRIVATE_AUTHENTICATION, AZURE_SAS_TOKEN
	AuthenticationType *SourceCdnRespAuthenticationType `json:"authentication_type,omitempty"`
}

func (o SourceCdnResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceCdnResp struct{}"
	}

	return strings.Join([]string{"SourceCdnResp", string(data)}, " ")
}

type SourceCdnRespProtocol struct {
	value string
}

type SourceCdnRespProtocolEnum struct {
	HTTP  SourceCdnRespProtocol
	HTTPS SourceCdnRespProtocol
}

func GetSourceCdnRespProtocolEnum() SourceCdnRespProtocolEnum {
	return SourceCdnRespProtocolEnum{
		HTTP: SourceCdnRespProtocol{
			value: "http",
		},
		HTTPS: SourceCdnRespProtocol{
			value: "https",
		},
	}
}

func (c SourceCdnRespProtocol) Value() string {
	return c.value
}

func (c SourceCdnRespProtocol) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceCdnRespProtocol) UnmarshalJSON(b []byte) error {
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

type SourceCdnRespAuthenticationType struct {
	value string
}

type SourceCdnRespAuthenticationTypeEnum struct {
	NONE                         SourceCdnRespAuthenticationType
	QINIU_PRIVATE_AUTHENTICATION SourceCdnRespAuthenticationType
	ALIYUN_OSS_A                 SourceCdnRespAuthenticationType
	ALIYUN_OSS_B                 SourceCdnRespAuthenticationType
	ALIYUN_OSS_C                 SourceCdnRespAuthenticationType
	KSYUN_PRIVATE_AUTHENTICATION SourceCdnRespAuthenticationType
	AZURE_SAS_TOKEN              SourceCdnRespAuthenticationType
}

func GetSourceCdnRespAuthenticationTypeEnum() SourceCdnRespAuthenticationTypeEnum {
	return SourceCdnRespAuthenticationTypeEnum{
		NONE: SourceCdnRespAuthenticationType{
			value: "NONE",
		},
		QINIU_PRIVATE_AUTHENTICATION: SourceCdnRespAuthenticationType{
			value: "QINIU_PRIVATE_AUTHENTICATION",
		},
		ALIYUN_OSS_A: SourceCdnRespAuthenticationType{
			value: "ALIYUN_OSS_A",
		},
		ALIYUN_OSS_B: SourceCdnRespAuthenticationType{
			value: "ALIYUN_OSS_B",
		},
		ALIYUN_OSS_C: SourceCdnRespAuthenticationType{
			value: "ALIYUN_OSS_C",
		},
		KSYUN_PRIVATE_AUTHENTICATION: SourceCdnRespAuthenticationType{
			value: "KSYUN_PRIVATE_AUTHENTICATION",
		},
		AZURE_SAS_TOKEN: SourceCdnRespAuthenticationType{
			value: "AZURE_SAS_TOKEN",
		},
	}
}

func (c SourceCdnRespAuthenticationType) Value() string {
	return c.value
}

func (c SourceCdnRespAuthenticationType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceCdnRespAuthenticationType) UnmarshalJSON(b []byte) error {
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
