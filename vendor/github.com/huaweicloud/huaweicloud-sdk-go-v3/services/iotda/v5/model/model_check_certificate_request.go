package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CheckCertificateRequest struct {

	// **参数说明**：实例ID。物理多租下各实例的唯一标识，一般华为云租户无需携带该参数，仅在物理多租场景下从管理面访问API时需要携带该参数。
	InstanceId *string `json:"Instance-Id,omitempty"`

	// **参数说明**：设备CA证书ID，在上传设备CA证书时由平台分配的唯一标识。
	CertificateId string `json:"certificate_id"`

	// **参数说明**：对证书执行的操作。 **取值范围**：当前仅支持verify:校验证书。
	ActionId string `json:"action_id"`

	Body *VerifyCertificateDto `json:"body,omitempty"`
}

func (o CheckCertificateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckCertificateRequest struct{}"
	}

	return strings.Join([]string{"CheckCertificateRequest", string(data)}, " ")
}
