package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 容器实例名称
type ContainerName struct {
}

func (o ContainerName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContainerName struct{}"
	}

	return strings.Join([]string{"ContainerName", string(data)}, " ")
}
