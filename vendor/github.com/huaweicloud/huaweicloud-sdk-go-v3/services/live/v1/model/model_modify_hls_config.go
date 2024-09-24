package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyHlsConfig struct {

	// 推流域名
	PushDomain string `json:"push_domain"`

	// 推流域名APP配置
	Application []PushDomainApplication `json:"application"`
}

func (o ModifyHlsConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyHlsConfig struct{}"
	}

	return strings.Join([]string{"ModifyHlsConfig", string(data)}, " ")
}
