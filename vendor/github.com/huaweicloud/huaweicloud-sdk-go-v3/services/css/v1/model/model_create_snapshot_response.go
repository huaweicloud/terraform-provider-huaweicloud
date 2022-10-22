package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateSnapshotResponse struct {
	Backup         *BackupRsp `json:"backup,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o CreateSnapshotResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSnapshotResponse struct{}"
	}

	return strings.Join([]string{"CreateSnapshotResponse", string(data)}, " ")
}
