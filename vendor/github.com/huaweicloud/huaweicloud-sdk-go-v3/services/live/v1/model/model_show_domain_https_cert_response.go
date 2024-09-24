package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowDomainHttpsCertResponse Response Object
type ShowDomainHttpsCertResponse struct {

	// 证书格式，默认为PEM，当前只支持PEM格式
	CertificateFormat *ShowDomainHttpsCertResponseCertificateFormat `json:"certificate_format,omitempty"`

	// 证书内容
	Certificate *string `json:"certificate,omitempty"`

	// 私钥内容
	CertificateKey *string `json:"certificate_key,omitempty"`

	// 是否开启重定向，默认false
	ForceRedirect *bool `json:"force_redirect,omitempty"`

	GmCertificate *GmCertificateInfo `json:"gm_certificate,omitempty"`

	TlsCertificate *TlsCertificateInfo `json:"tls_certificate,omitempty"`
	HttpStatusCode int                 `json:"-"`
}

func (o ShowDomainHttpsCertResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainHttpsCertResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainHttpsCertResponse", string(data)}, " ")
}

type ShowDomainHttpsCertResponseCertificateFormat struct {
	value string
}

type ShowDomainHttpsCertResponseCertificateFormatEnum struct {
	PEM ShowDomainHttpsCertResponseCertificateFormat
}

func GetShowDomainHttpsCertResponseCertificateFormatEnum() ShowDomainHttpsCertResponseCertificateFormatEnum {
	return ShowDomainHttpsCertResponseCertificateFormatEnum{
		PEM: ShowDomainHttpsCertResponseCertificateFormat{
			value: "PEM",
		},
	}
}

func (c ShowDomainHttpsCertResponseCertificateFormat) Value() string {
	return c.value
}

func (c ShowDomainHttpsCertResponseCertificateFormat) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowDomainHttpsCertResponseCertificateFormat) UnmarshalJSON(b []byte) error {
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
