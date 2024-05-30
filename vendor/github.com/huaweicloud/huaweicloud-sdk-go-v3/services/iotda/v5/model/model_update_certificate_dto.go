package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateCertificateDto 更新CA证书结构体。
type UpdateCertificateDto struct {

	// 是否开启自注册能力，当为true时该功能必须配合预调配功能使用，true：是，false：否。
	ProvisionEnable *bool `json:"provision_enable,omitempty"`

	// 预调配模板ID，该CA证书绑定的预调配模板id，当该字段传null时表示解除绑定关系。
	TemplateId *string `json:"template_id,omitempty"`
}

func (o UpdateCertificateDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateCertificateDto struct{}"
	}

	return strings.Join([]string{"UpdateCertificateDto", string(data)}, " ")
}
