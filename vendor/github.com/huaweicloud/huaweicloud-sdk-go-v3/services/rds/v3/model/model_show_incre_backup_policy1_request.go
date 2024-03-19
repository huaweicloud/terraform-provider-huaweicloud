package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIncreBackupPolicy1Request Request Object
type ShowIncreBackupPolicy1Request struct {

	// instance id
	InstanceId string `json:"instance_id"`
}

func (o ShowIncreBackupPolicy1Request) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIncreBackupPolicy1Request struct{}"
	}

	return strings.Join([]string{"ShowIncreBackupPolicy1Request", string(data)}, " ")
}
