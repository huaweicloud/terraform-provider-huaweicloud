package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowLogBackupRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`

	Body *GetLogBackupReq `json:"body,omitempty"`
}

func (o ShowLogBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowLogBackupRequest struct{}"
	}

	return strings.Join([]string{"ShowLogBackupRequest", string(data)}, " ")
}
