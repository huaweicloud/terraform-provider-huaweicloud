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
}

func (o CreateCnfReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCnfReq struct{}"
	}

	return strings.Join([]string{"CreateCnfReq", string(data)}, " ")
}
