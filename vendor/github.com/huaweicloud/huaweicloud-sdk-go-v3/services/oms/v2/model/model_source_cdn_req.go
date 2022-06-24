package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源端CDN配置。
type SourceCdnReq struct {

	// CDN鉴权秘钥，如果CDN需要进行鉴权，则此选项为必选。  无需授权：无需配置此项。 Qiniu：无需配置此项。 Aliyun：根据authentication_type指定的鉴权方式配置此项。 KingsoftCloud：无需配置此项。
	AuthenticationKey *string `json:"authentication_key,omitempty"`

	// 鉴权类型: NONE, QINIU_PRIVATE_AUTHENTICATION, ALIYUN_OSS_A, ALIYUN_OSS_B, ALIYUN_OSS_C, KSYUN_PRIVATE_AUTHENTICATION, AZURE_SAS_TOKEN
	AuthenticationType SourceCdnReqAuthenticationType `json:"authentication_type"`

	//   从指定域名获取对象。
	Domain string `json:"domain"`

	// 协议类型，支持http和https协议。
	Protocol SourceCdnReqProtocol `json:"protocol"`
}

func (o SourceCdnReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SourceCdnReq struct{}"
	}

	return strings.Join([]string{"SourceCdnReq", string(data)}, " ")
}

type SourceCdnReqAuthenticationType struct {
	value string
}

type SourceCdnReqAuthenticationTypeEnum struct {
	NONE                         SourceCdnReqAuthenticationType
	QINIU_PRIVATE_AUTHENTICATION SourceCdnReqAuthenticationType
	ALIYUN_OSS_A                 SourceCdnReqAuthenticationType
	ALIYUN_OSS_B                 SourceCdnReqAuthenticationType
	ALIYUN_OSS_C                 SourceCdnReqAuthenticationType
	KSYUN_PRIVATE_AUTHENTICATION SourceCdnReqAuthenticationType
}

func GetSourceCdnReqAuthenticationTypeEnum() SourceCdnReqAuthenticationTypeEnum {
	return SourceCdnReqAuthenticationTypeEnum{
		NONE: SourceCdnReqAuthenticationType{
			value: "NONE",
		},
		QINIU_PRIVATE_AUTHENTICATION: SourceCdnReqAuthenticationType{
			value: "QINIU_PRIVATE_AUTHENTICATION",
		},
		ALIYUN_OSS_A: SourceCdnReqAuthenticationType{
			value: "ALIYUN_OSS_A",
		},
		ALIYUN_OSS_B: SourceCdnReqAuthenticationType{
			value: "ALIYUN_OSS_B",
		},
		ALIYUN_OSS_C: SourceCdnReqAuthenticationType{
			value: "ALIYUN_OSS_C",
		},
		KSYUN_PRIVATE_AUTHENTICATION: SourceCdnReqAuthenticationType{
			value: "KSYUN_PRIVATE_AUTHENTICATION",
		},
	}
}

func (c SourceCdnReqAuthenticationType) Value() string {
	return c.value
}

func (c SourceCdnReqAuthenticationType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceCdnReqAuthenticationType) UnmarshalJSON(b []byte) error {
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

type SourceCdnReqProtocol struct {
	value string
}

type SourceCdnReqProtocolEnum struct {
	HTTP  SourceCdnReqProtocol
	HTTPS SourceCdnReqProtocol
}

func GetSourceCdnReqProtocolEnum() SourceCdnReqProtocolEnum {
	return SourceCdnReqProtocolEnum{
		HTTP: SourceCdnReqProtocol{
			value: "http",
		},
		HTTPS: SourceCdnReqProtocol{
			value: "https",
		},
	}
}

func (c SourceCdnReqProtocol) Value() string {
	return c.value
}

func (c SourceCdnReqProtocol) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SourceCdnReqProtocol) UnmarshalJSON(b []byte) error {
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
