package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClustersDetailsRequest Request Object
type ListClustersDetailsRequest struct {

	// 指定查询起始值，默认值为1，即从第1个集群开始查询。
	Start *int32 `json:"start,omitempty"`

	// 指定查询个数，默认值为10，即一次查询10个集群信息。
	Limit *int32 `json:"limit,omitempty"`

	// 指定查询的集群引擎类型。
	DatastoreType *string `json:"datastoreType,omitempty"`
}

func (o ListClustersDetailsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClustersDetailsRequest struct{}"
	}

	return strings.Join([]string{"ListClustersDetailsRequest", string(data)}, " ")
}
