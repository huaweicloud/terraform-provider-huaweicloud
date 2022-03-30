package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlProxyFlavorsRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`
}

func (o ShowGaussMySqlProxyFlavorsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlProxyFlavorsRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlProxyFlavorsRequest", string(data)}, " ")
}
