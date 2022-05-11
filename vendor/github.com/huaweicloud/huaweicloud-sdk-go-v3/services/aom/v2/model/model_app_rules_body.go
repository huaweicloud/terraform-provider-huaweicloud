package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AppRulesBody struct {

	// 服务参数。
	AppRules *[]AppRules `json:"appRules,omitempty"`
}

func (o AppRulesBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppRulesBody struct{}"
	}

	return strings.Join([]string{"AppRulesBody", string(data)}, " ")
}
