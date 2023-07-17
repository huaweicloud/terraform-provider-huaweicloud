package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateOtaPackage 添加升级包关联到OBS对象结构体
type CreateOtaPackage struct {

	// **参数说明**：资源空间ID。存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的升级包归属到哪个资源空间下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId string `json:"app_id"`

	// **参数说明**：升级包类型。 **取值范围**：软件包必须设置为：softwarePackage，固件包必须设置为：firmwarePackage。
	PackageType string `json:"package_type"`

	// **参数说明**：设备关联的产品ID，用于唯一标识一个产品模型，创建产品后获得。方法请参见 [[创建产品](https://support.huaweicloud.com/api-iothub/iot_06_v5_0050.html)](tag:hws)[[创建产品](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0050.html)](tag:hws_hk)。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	ProductId string `json:"product_id"`

	// **参数说明**：升级包版本号。 **取值范围**：长度不超过256，只允许字母、数字、下划线（_）、连接符（-）、英文点（.）的组合。
	Version string `json:"version"`

	// **参数说明**：支持用于升级此版本包的设备源版本号列表。最多支持20个源版本号。 **取值范围**：源版本号列表，源版本号只允许字母、数字、下划线（_）、连接符（-）、英文点（.）的组合。
	SupportSourceVersions *[]string `json:"support_source_versions,omitempty"`

	// **参数说明**：用于描述升级包的功能等信息。 **取值范围**：长度不超过1024。
	Description *string `json:"description,omitempty"`

	// **参数说明**：推送给设备的自定义信息。添加该升级包完成，并创建升级任务后，物联网平台向设备下发升级通知时，会下发该自定义信息给设备。 **取值范围**：长度不超过4096。
	CustomInfo *string `json:"custom_info,omitempty"`

	FileLocation *FileLocation `json:"file_location"`
}

func (o CreateOtaPackage) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOtaPackage struct{}"
	}

	return strings.Join([]string{"CreateOtaPackage", string(data)}, " ")
}
