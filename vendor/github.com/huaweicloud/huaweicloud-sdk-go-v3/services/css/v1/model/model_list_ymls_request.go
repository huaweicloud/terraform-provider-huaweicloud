package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListYmlsRequest struct {

	// 指定查询集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListYmlsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListYmlsRequest struct{}"
	}

	return strings.Join([]string{"ListYmlsRequest", string(data)}, " ")
}
