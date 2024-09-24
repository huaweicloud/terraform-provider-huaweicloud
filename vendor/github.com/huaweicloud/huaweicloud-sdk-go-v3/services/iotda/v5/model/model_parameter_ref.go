package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ParameterRef struct {

	// 参数引用名称
	Ref string `json:"ref"`
}

func (o ParameterRef) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ParameterRef struct{}"
	}

	return strings.Join([]string{"ParameterRef", string(data)}, " ")
}
