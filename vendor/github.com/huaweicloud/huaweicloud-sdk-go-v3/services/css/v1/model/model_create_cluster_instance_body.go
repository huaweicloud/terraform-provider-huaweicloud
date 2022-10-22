package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 实例对象。
type CreateClusterInstanceBody struct {

	// 实例规格名称。可以使用[获取实例规格列表](ListFlavors.xml)的name属性确认当前拥有的规格信息。
	FlavorRef string `json:"flavorRef"`

	Volume *CreateClusterInstanceVolumeBody `json:"volume"`

	Nics *CreateClusterInstanceNicsBody `json:"nics"`

	// 可用区。不填时默认创建单AZ。  如果需要创建多AZ，各个AZ之间使用英文逗号分隔，比如az1,az2 ，az不能重复输入，并且要求节点个数大于等于AZ个数。  如果节点个数为AZ个数的倍数，节点将会均匀的分布到各个AZ。如果节点个数不为AZ个数的倍数时，各个AZ分布的节点数量之差的绝对值小于等于1。
	AvailabilityZone *string `json:"availability_zone,omitempty"`
}

func (o CreateClusterInstanceBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterInstanceBody struct{}"
	}

	return strings.Join([]string{"CreateClusterInstanceBody", string(data)}, " ")
}
