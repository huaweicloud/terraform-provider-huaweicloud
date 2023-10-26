package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterInstanceBody 实例对象。
type CreateClusterInstanceBody struct {

	// 实例规格名称。可以使用[获取实例规格列表](ListFlavors.xml)的name属性确认当前拥有的规格信息。
	FlavorRef string `json:"flavorRef"`

	Volume *CreateClusterInstanceVolumeBody `json:"volume"`

	Nics *CreateClusterInstanceNicsBody `json:"nics"`

	// 可用区。需要指定可用区的名称（可用分区名称）。 默认指定单AZ。指定多AZ时，各个可用分区名称需要使用英文逗号（,）分隔，以“华北-北京四”为例，选择三AZ时，availability_zone取值为cn-north-4a,cn-north-4b,cn-north-4c。如果使用单AZ，availability_zone默认取值为空。 >说明   选择多AZ时，各个可用分区名称不能重复输入，并且要求节点个数大于等于AZ个数。      如果节点个数为AZ个数的倍数，节点将会均匀的分布到各个AZ。如果节点个数不为AZ个数的倍数时，各个AZ分布的节点数量之差的绝对值小于等于1。     可用分区名称，请在[[地区和终端节点](https://developer.huaweicloud.com/endpoint?CSS)](tag:hws)[[地区和终端节点](https://developer.huaweicloud.com/intl/zh-cn/endpoint?CSS)](tag:hk_hws)获取。
	AvailabilityZone *string `json:"availability_zone,omitempty"`
}

func (o CreateClusterInstanceBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterInstanceBody struct{}"
	}

	return strings.Join([]string{"CreateClusterInstanceBody", string(data)}, " ")
}
