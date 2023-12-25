package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddonInstanceRollbackRequest struct {

	// 集群ID
	ClusterID string `json:"clusterID"`
}

func (o AddonInstanceRollbackRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddonInstanceRollbackRequest struct{}"
	}

	return strings.Join([]string{"AddonInstanceRollbackRequest", string(data)}, " ")
}
