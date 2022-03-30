package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateGaussMySqlBackupResponse struct {
	Backup *Backup `json:"backup,omitempty"`
	// 任务ID。

	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateGaussMySqlBackupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateGaussMySqlBackupResponse struct{}"
	}

	return strings.Join([]string{"CreateGaussMySqlBackupResponse", string(data)}, " ")
}
