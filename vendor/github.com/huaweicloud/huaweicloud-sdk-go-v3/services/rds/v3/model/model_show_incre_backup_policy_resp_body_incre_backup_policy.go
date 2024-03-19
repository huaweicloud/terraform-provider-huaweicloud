package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIncreBackupPolicyRespBodyIncreBackupPolicy incre backup policy
type ShowIncreBackupPolicyRespBodyIncreBackupPolicy struct {

	// 增备时间间隔（分）
	Interval *int32 `json:"interval,omitempty"`
}

func (o ShowIncreBackupPolicyRespBodyIncreBackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIncreBackupPolicyRespBodyIncreBackupPolicy struct{}"
	}

	return strings.Join([]string{"ShowIncreBackupPolicyRespBodyIncreBackupPolicy", string(data)}, " ")
}
