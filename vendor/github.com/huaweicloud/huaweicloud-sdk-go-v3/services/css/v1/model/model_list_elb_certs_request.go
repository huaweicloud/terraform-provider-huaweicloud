package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListElbCertsRequest Request Object
type ListElbCertsRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ListElbCertsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListElbCertsRequest struct{}"
	}

	return strings.Join([]string{"ListElbCertsRequest", string(data)}, " ")
}
