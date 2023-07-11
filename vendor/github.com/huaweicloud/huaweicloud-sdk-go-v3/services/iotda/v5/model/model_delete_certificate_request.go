package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteCertificateRequest Request Object
type DeleteCertificateRequest struct {

	// Sp用户Token。通过调用IoBPS服务获取SP用户Token
	SpAuthToken *string `json:"Sp-Auth-Token,omitempty"`

	// Stage用户的Token, 仅提供给IoStage服务使用
	StageAuthToken *string `json:"Stage-Auth-Token,omitempty"`

	// 实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// 设备CA证书ID，在上传设备CA证书时由平台分配的唯一标识。
	CertificateId string `json:"certificate_id"`
}

func (o DeleteCertificateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteCertificateRequest struct{}"
	}

	return strings.Join([]string{"DeleteCertificateRequest", string(data)}, " ")
}
