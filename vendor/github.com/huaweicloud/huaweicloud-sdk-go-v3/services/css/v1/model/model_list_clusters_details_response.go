package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListClustersDetailsResponse struct {

	// 集群对象列表。
	Clusters       *[]ClusterList `json:"clusters,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListClustersDetailsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClustersDetailsResponse struct{}"
	}

	return strings.Join([]string{"ListClustersDetailsResponse", string(data)}, " ")
}
