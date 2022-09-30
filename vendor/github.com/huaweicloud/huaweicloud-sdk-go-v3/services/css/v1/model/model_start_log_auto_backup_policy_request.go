package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type StartLogAutoBackupPolicyRequest struct {

	// 指定开启日志备份的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *StartLogAutoBackupPolicyReq `json:"body,omitempty"`
}

func (o StartLogAutoBackupPolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogAutoBackupPolicyRequest struct{}"
	}

	return strings.Join([]string{"StartLogAutoBackupPolicyRequest", string(data)}, " ")
}
