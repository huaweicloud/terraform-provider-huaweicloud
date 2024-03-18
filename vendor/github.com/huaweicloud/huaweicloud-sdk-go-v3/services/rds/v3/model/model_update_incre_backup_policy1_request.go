package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateIncreBackupPolicy1Request Request Object
type UpdateIncreBackupPolicy1Request struct {

	// instance id
	InstanceId string `json:"instance_id"`

	Body *UpdateIncreBackupPolicy1RequestBody `json:"body,omitempty"`
}

func (o UpdateIncreBackupPolicy1Request) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateIncreBackupPolicy1Request struct{}"
	}

	return strings.Join([]string{"UpdateIncreBackupPolicy1Request", string(data)}, " ")
}
