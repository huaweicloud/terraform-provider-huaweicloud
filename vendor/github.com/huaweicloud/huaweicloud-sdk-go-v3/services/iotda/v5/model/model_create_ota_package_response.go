package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateOtaPackageResponse Response Object
type CreateOtaPackageResponse struct {

	// **参数说明**：升级包ID，用于唯一标识一个升级包。由物联网平台分配获得。 **取值范围**：长度不超过36，只允许字母、数字、连接符（-）的组合。
	PackageId *string `json:"package_id,omitempty"`

	// **参数说明**：资源空间ID。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：升级包类型。 **取值范围**：软件包必须设置为：softwarePackage，固件包必须设置为：firmwarePackage。
	PackageType *string `json:"package_type,omitempty"`

	// **参数说明**：设备关联的产品ID，用于唯一标识一个产品模型，创建产品后获得。方法请参见 [[创建产品](https://support.huaweicloud.com/api-iothub/iot_06_v5_0050.html)](tag:hws)[[创建产品](https://support.huaweicloud.com/intl/zh-cn/api-iothub/iot_06_v5_0050.html)](tag:hws_hk)。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	ProductId *string `json:"product_id,omitempty"`

	// **参数说明**：升级包版本号。 **取值范围**：长度不超过256，只允许字母、数字、下划线（_）、连接符（-）、英文点（.）的组合。
	Version *string `json:"version,omitempty"`

	// **参数说明**：支持用于升级此版本包的设备源版本号列表。最多支持20个源版本号。 **取值范围**：源版本号列表，源版本号只允许字母、数字、下划线（_）、连接符（-）、英文点（.）的组合。
	SupportSourceVersions *[]string `json:"support_source_versions,omitempty"`

	// **参数说明**：用于描述升级包的功能等信息。 **取值范围**：长度不超过1024。
	Description *string `json:"description,omitempty"`

	// **参数说明**：推送给设备的自定义信息。添加该升级包完成，并创建升级任务后，物联网平台向设备下发升级通知时，会下发该自定义信息给设备。 **取值范围**：长度不超过4096。
	CustomInfo *string `json:"custom_info,omitempty"`

	// 软固件包上传到物联网平台的时间，格式：\"yyyyMMdd'T'HHmmss'Z'\"。
	CreateTime *string `json:"create_time,omitempty"`

	FileLocation   *FileLocation `json:"file_location,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o CreateOtaPackageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOtaPackageResponse struct{}"
	}

	return strings.Join([]string{"CreateOtaPackageResponse", string(data)}, " ")
}
