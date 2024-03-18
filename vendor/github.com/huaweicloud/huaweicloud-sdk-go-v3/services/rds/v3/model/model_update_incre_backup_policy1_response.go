package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateIncreBackupPolicy1Response Response Object
type UpdateIncreBackupPolicy1Response struct {

	// job id
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateIncreBackupPolicy1Response) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIncreBackupPolicy1Response struct{}"
	}

	return strings.Join([]string{"UpdateIncreBackupPolicy1Response", string(data)}, " ")
}
