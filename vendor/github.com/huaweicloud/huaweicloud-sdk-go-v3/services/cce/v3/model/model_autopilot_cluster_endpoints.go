package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AutopilotClusterEndpoints struct {

	// 集群中 kube-apiserver 的访问地址
	Url *string `json:"url,omitempty"`

	// 集群访问地址的类型 - Internal：用户子网内访问的地址 - External：公网访问的地址
	Type *string `json:"type,omitempty"`
}

func (o AutopilotClusterEndpoints) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterEndpoints struct{}"
	}

	return strings.Join([]string{"AutopilotClusterEndpoints", string(data)}, " ")
}
