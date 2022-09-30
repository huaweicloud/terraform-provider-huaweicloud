package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListYmlsJobRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListYmlsJobRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListYmlsJobRequest struct{}"
	}

	return strings.Join([]string{"ListYmlsJobRequest", string(data)}, " ")
}
