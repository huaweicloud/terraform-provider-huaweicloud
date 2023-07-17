package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateSnapshotSettingResponse Response Object
type UpdateSnapshotSettingResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateSnapshotSettingResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSnapshotSettingResponse struct{}"
	}

	return strings.Join([]string{"UpdateSnapshotSettingResponse", string(data)}, " ")
}
