package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StopBackupResponse Response Object
type StopBackupResponse struct {

	// 任务id
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StopBackupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StopBackupResponse struct{}"
	}

	return strings.Join([]string{"StopBackupResponse", string(data)}, " ")
}
