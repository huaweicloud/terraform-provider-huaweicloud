package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListContainerNodesResponse Response Object
type ListContainerNodesResponse struct {

	// 容器节点总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 容器节点列表
	DataList       *[]ContainerNodeInfo `json:"data_list,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ListContainerNodesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListContainerNodesResponse struct{}"
	}

	return strings.Join([]string{"ListContainerNodesResponse", string(data)}, " ")
}
