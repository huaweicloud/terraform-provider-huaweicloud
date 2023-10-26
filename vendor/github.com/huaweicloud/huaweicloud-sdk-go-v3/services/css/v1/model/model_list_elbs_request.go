package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListElbsRequest Request Object
type ListElbsRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListElbsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListElbsRequest struct{}"
	}

	return strings.Join([]string{"ListElbsRequest", string(data)}, " ")
}
