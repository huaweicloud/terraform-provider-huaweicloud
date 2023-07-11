package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowDomainKeyChainResponse Response Object
type ShowDomainKeyChainResponse struct {

	// 防盗链Key值，由32个字符组成，支持大写字母、小写字母、数字。不可为纯数字或纯字母。
	Key *string `json:"key,omitempty"`

	// 计算鉴权串的方式： - d_sha256：鉴权方式D，采用HMAC-SHA256算法，建议优先选择此方式； - c_aes：鉴权方式C，采用对称加密算法； - b_md5：鉴权方式B，采用MD5信息摘要算法； - a_md5：鉴权方式A，采用MD5信息摘要算法。  > 鉴权方式ABC存在安全风险，鉴权方式D拥有更高的安全性，建议您优先使用鉴权方式D。
	AuthType *ShowDomainKeyChainResponseAuthType `json:"auth_type,omitempty"`

	// URL鉴权信息的超时时长  取值范围：[60，2592000]，即1分钟-30天  单位：秒  鉴权信息中携带的请求时间与直播服务收到请求时的时间的最大差值，用于检查直播推流URL或者直播播放URL是否已过期
	Timeout        *int32 `json:"timeout,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowDomainKeyChainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainKeyChainResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainKeyChainResponse", string(data)}, " ")
}

type ShowDomainKeyChainResponseAuthType struct {
	value string
}

type ShowDomainKeyChainResponseAuthTypeEnum struct {
	D_SHA256 ShowDomainKeyChainResponseAuthType
	C_AES    ShowDomainKeyChainResponseAuthType
	B_MD5    ShowDomainKeyChainResponseAuthType
	A_MD5    ShowDomainKeyChainResponseAuthType
}

func GetShowDomainKeyChainResponseAuthTypeEnum() ShowDomainKeyChainResponseAuthTypeEnum {
	return ShowDomainKeyChainResponseAuthTypeEnum{
		D_SHA256: ShowDomainKeyChainResponseAuthType{
			value: "d_sha256",
		},
		C_AES: ShowDomainKeyChainResponseAuthType{
			value: "c_aes",
		},
		B_MD5: ShowDomainKeyChainResponseAuthType{
			value: "b_md5",
		},
		A_MD5: ShowDomainKeyChainResponseAuthType{
			value: "a_md5",
		},
	}
}

func (c ShowDomainKeyChainResponseAuthType) Value() string {
	return c.value
}

func (c ShowDomainKeyChainResponseAuthType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowDomainKeyChainResponseAuthType) UnmarshalJSON(b []byte) error {
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
