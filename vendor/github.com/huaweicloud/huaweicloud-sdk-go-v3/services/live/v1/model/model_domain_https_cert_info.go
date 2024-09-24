package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type DomainHttpsCertInfo struct {

	// 证书格式，默认为PEM，当前只支持PEM格式
	CertificateFormat *DomainHttpsCertInfoCertificateFormat `json:"certificate_format,omitempty"`

	// 证书内容
	Certificate *string `json:"certificate,omitempty"`

	// 私钥内容
	CertificateKey *string `json:"certificate_key,omitempty"`

	// 是否开启重定向，默认false
	ForceRedirect *bool `json:"force_redirect,omitempty"`

	GmCertificate *GmCertificateInfo `json:"gm_certificate,omitempty"`

	TlsCertificate *TlsCertificateInfo `json:"tls_certificate,omitempty"`
}

func (o DomainHttpsCertInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainHttpsCertInfo struct{}"
	}

	return strings.Join([]string{"DomainHttpsCertInfo", string(data)}, " ")
}

type DomainHttpsCertInfoCertificateFormat struct {
	value string
}

type DomainHttpsCertInfoCertificateFormatEnum struct {
	PEM DomainHttpsCertInfoCertificateFormat
}

func GetDomainHttpsCertInfoCertificateFormatEnum() DomainHttpsCertInfoCertificateFormatEnum {
	return DomainHttpsCertInfoCertificateFormatEnum{
		PEM: DomainHttpsCertInfoCertificateFormat{
			value: "PEM",
		},
	}
}

func (c DomainHttpsCertInfoCertificateFormat) Value() string {
	return c.value
}

func (c DomainHttpsCertInfoCertificateFormat) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DomainHttpsCertInfoCertificateFormat) UnmarshalJSON(b []byte) error {
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
