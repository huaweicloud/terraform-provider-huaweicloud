package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 参数模板信息。
type ConfigurationSummary struct {
	// 参数组ID。

	Id string `json:"id"`
	// 参数组名称。

	Name string `json:"name"`
	// 参数组描述。

	Description *string `json:"description,omitempty"`
	// 引擎版本。

	DatastoreVersionName string `json:"datastore_version_name"`
	// 引擎名。

	DatastoreName string `json:"datastore_name"`
	// 创建时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	Created string `json:"created"`
	// 更新时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	Updated string `json:"updated"`
	// 是否是用户自定义参数模板：  - false，表示为系统默认参数模板。 - true，表示为用户自定义参数模板。

	UserDefined bool `json:"user_defined"`
}

func (o ConfigurationSummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationSummary struct{}"
	}

	return strings.Join([]string{"ConfigurationSummary", string(data)}, " ")
}
