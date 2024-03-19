package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopBackupRequest Request Object
type StopBackupRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o StopBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopBackupRequest struct{}"
	}

	return strings.Join([]string{"StopBackupRequest", string(data)}, " ")
}
