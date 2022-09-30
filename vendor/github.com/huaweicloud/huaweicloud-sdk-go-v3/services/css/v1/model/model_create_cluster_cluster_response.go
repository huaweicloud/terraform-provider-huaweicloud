package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 集群对象。
type CreateClusterClusterResponse struct {

	// 集群ID。
	Id *string `json:"id,omitempty"`

	// 集群名称。
	Name *string `json:"name,omitempty"`
}

func (o CreateClusterClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterClusterResponse struct{}"
	}

	return strings.Join([]string{"CreateClusterClusterResponse", string(data)}, " ")
}
