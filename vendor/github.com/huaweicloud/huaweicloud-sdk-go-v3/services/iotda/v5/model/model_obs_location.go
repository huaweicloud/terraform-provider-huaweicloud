package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ObsLocation 升级包在OBS的位置
type ObsLocation struct {

	// **参数说明**：OBS所在区域。您可以从[[地区和终端节点](https://developer.huaweicloud.com/endpoint?OBS)](tag:hws)[[地区和终端节点](https://developer.huaweicloud.com/intl/zh-cn/endpoint?OBS)](tag:hws_hk)中查询服务的终端节点。 **取值范围**：长度不超过256，只允许字母、数字、连接符（-）的组合。
	RegionName string `json:"region_name"`

	// **参数说明**：OBS桶名称。 **取值范围**：长度最小为3，最大为63，只允许小写字母、数字、连接符（-）、英文点（.）的组合。
	BucketName string `json:"bucket_name"`

	// **参数说明**：OBS对象名称(包含文件夹路径)。 **取值范围**：长度不超过1024。
	ObjectKey string `json:"object_key"`

	// **参数说明**：SHA256算法计算出的升级包签名值。添加该升级包完成，并创建升级任务后，物联网平台向设备下发升级通知时，会下发该签名给设备。 **取值范围**：长度为64，只允许大小写字母a到f、数字的组合。
	Sign *string `json:"sign,omitempty"`
}

func (o ObsLocation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsLocation struct{}"
	}

	return strings.Join([]string{"ObsLocation", string(data)}, " ")
}
