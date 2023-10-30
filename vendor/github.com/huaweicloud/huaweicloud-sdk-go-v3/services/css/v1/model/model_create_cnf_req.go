package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateCnfReq struct {

	// 配置文件名称。4～32个字符，只能包含数字、字母、中划线和下划线，且必须以字母开头。
	Name string `json:"name"`

	// 配置文件内容。
	ConfContent string `json:"confContent"`

	Setting *Setting `json:"setting"`

	// 敏感字符替换 输入需要隐藏的敏感字串列表。配置隐藏字符串列表后，在返回的配置内容中，会将所有在列表中的字串隐藏为***（列表最大支持20条，单个字串最大长度512字节）
	SensitiveWords *[]string `json:"sensitive_words,omitempty"`
}

func (o CreateCnfReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCnfReq struct{}"
	}

	return strings.Join([]string{"CreateCnfReq", string(data)}, " ")
}
