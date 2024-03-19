package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CertificatesRspDto struct {

	// CA证书ID，在上传CA证书时由平台分配的唯一标识。
	CertificateId *string `json:"certificate_id,omitempty"`

	// CA证书CN名称。
	CnName *string `json:"cn_name,omitempty"`

	// CA证书所有者。
	Owner *string `json:"owner,omitempty"`

	// CA证书验证状态。true代表证书已通过验证，可进行设备证书认证接入。false代表证书未通过验证。
	Status *bool `json:"status,omitempty"`

	// CA证书验证码。
	VerifyCode *string `json:"verify_code,omitempty"`

	// 是否开启自注册能力，当为true时该功能必须配合自注册模板使用，true：是，false：否。
	ProvisionEnable *bool `json:"provision_enable,omitempty"`

	// 绑定的自注册模板ID。
	TemplateId *string `json:"template_id,omitempty"`

	// 创建证书日期。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateDate *string `json:"create_date,omitempty"`

	// CA证书生效日期。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	EffectiveDate *string `json:"effective_date,omitempty"`

	// CA证书失效日期。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	ExpiryDate *string `json:"expiry_date,omitempty"`
}

func (o CertificatesRspDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CertificatesRspDto struct{}"
	}

	return strings.Join([]string{"CertificatesRspDto", string(data)}, " ")
}
