package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Statement 策略文档结构。
type Statement struct {

	// 指定是允许还是拒绝该操作。既有允许（ALLOW）又有拒绝（DENY）的授权语句时，遵循拒绝（DENY）优先的原则。 - ALLOW：允许。 - DENY：拒绝。
	Effect string `json:"effect"`

	// 用于指定策略允许或拒绝的操作。格式为：服务名:资源:操作。当前支持的操作类型如下： - iotda:devices:publish：设备使用MQTT协议发布消息。 - iotda:devices:subscribe：设备使用MQTT协议订阅消息。
	Actions []string `json:"actions"`

	// 用于指定允许或拒绝对其执行操作的资源。格式为：资源类型:资源名称。如设备订阅的资源为：topic:/v1/${devices.deviceId}/test/hello。  **取值范围**：资源列表长度最小为1，最大为10，列表中的资源取值范围：仅支持字母，数字，以及/{}$=+#?*:._-组合。
	Resources []string `json:"resources"`
}

func (o Statement) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Statement struct{}"
	}

	return strings.Join([]string{"Statement", string(data)}, " ")
}
