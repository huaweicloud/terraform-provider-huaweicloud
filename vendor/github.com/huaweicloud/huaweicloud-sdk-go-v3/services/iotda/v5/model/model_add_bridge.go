package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddBridge 添加网桥结构体。
type AddBridge struct {

	// 网桥名称。**取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	BridgeName string `json:"bridge_name"`

	// 网桥ID。**取值范围**：长度不超过36，只允许字母、数字、_-字符的组合。
	BridgeId *string `json:"bridge_id,omitempty"`
}

func (o AddBridge) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddBridge struct{}"
	}

	return strings.Join([]string{"AddBridge", string(data)}, " ")
}
