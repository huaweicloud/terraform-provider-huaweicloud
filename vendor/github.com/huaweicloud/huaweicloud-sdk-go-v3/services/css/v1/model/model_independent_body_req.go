package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type IndependentBodyReq struct {

	// 规格id，该参数通过[获取实例规格列表](ListFlavors.xml)接口获取，根据集群版本选择所需要的规格id
	FlavorRef string `json:"flavor_ref"`

	// 要独立节点个数。 - 如果路径参数type取值为“ess-master”即新增独立master节点，节点个数必须为大于等于三且小于等于10的奇数。 - 如果路径参数type取值为“ess-client”即新增独立client节点，节点个数要求大于等于1小于等于32。
	NodeSize int32 `json:"node_size"`

	// 节点存储类型：取值为ULTRAHIGH，COMMON，HIGH。
	VolumeType string `json:"volume_type"`
}

func (o IndependentBodyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IndependentBodyReq struct{}"
	}

	return strings.Join([]string{"IndependentBodyReq", string(data)}, " ")
}
