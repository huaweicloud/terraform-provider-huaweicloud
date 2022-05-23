package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 命令参数响应对象。
type ServiceCommandResponse struct {

	// **参数说明**：设备命令响应名称。 **取值范围**：长度不超过128，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	ResponseName string `json:"response_name"`

	// **参数说明**：设备命令响应的参数列表。
	Paras *[]ServiceCommandPara `json:"paras,omitempty"`
}

func (o ServiceCommandResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceCommandResponse struct{}"
	}

	return strings.Join([]string{"ServiceCommandResponse", string(data)}, " ")
}
