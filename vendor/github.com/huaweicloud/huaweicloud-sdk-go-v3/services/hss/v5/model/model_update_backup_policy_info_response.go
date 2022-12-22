package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateBackupPolicyInfoResponse struct {

	// 错误编码,成功返回0
	ErrorCode *int32 `json:"error_code,omitempty"`

	// 错误描述，成功返回success
	ErrorDescription *string `json:"error_description,omitempty"`
	HttpStatusCode   int     `json:"-"`
}

func (o UpdateBackupPolicyInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBackupPolicyInfoResponse struct{}"
	}

	return strings.Join([]string{"UpdateBackupPolicyInfoResponse", string(data)}, " ")
}
