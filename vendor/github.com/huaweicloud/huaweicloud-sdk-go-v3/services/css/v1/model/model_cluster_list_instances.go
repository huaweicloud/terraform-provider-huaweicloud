package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 节点对象。
type ClusterListInstances struct {

	// 节点状态值。  - 100：创建中。 - 200：可用。 - 303：不可用，如创建失败。
	Status *string `json:"status,omitempty"`

	// 当前节点的类型。
	Type *string `json:"type,omitempty"`

	// 实例ID。
	Id *string `json:"id,omitempty"`

	// 实例名字。
	Name *string `json:"name,omitempty"`

	// 节点规格名称。
	SpecCode *string `json:"specCode,omitempty"`

	// 节点所属AZ信息。
	AzCode *string `json:"azCode,omitempty"`

	// 实例ip信息。
	Ip *string `json:"ip,omitempty"`

	Volume *ClusterVolumeRsp `json:"volume,omitempty"`
}

func (o ClusterListInstances) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterListInstances struct{}"
	}

	return strings.Join([]string{"ClusterListInstances", string(data)}, " ")
}
