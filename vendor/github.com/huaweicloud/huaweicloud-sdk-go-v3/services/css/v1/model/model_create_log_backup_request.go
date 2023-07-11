package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateLogBackupRequest Request Object
type CreateLogBackupRequest struct {

	// 指定待备份日志的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o CreateLogBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLogBackupRequest struct{}"
	}

	return strings.Join([]string{"CreateLogBackupRequest", string(data)}, " ")
}
