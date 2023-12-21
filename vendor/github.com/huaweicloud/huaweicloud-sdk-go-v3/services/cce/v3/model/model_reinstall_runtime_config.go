package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ReinstallRuntimeConfig 节点重装场景容器运行时配置
type ReinstallRuntimeConfig struct {

	// 节点上单容器的可用磁盘空间大小，单位G。  不配置该值或值为0时将使用默认值，Devicemapper模式下默认值为10；OverlayFS模式默认不限制单容器可用空间大小，且dockerBaseSize设置仅在新版本集群的EulerOS节点上生效。  CCE节点容器运行时空间配置请参考[数据盘空间分配说明](cce_01_0341.xml)。  Devicemapper模式下建议dockerBaseSize配置不超过80G，设置过大时可能会导致容器运行时初始化时间过长而启动失败，若对容器磁盘大小有特殊要求，可考虑使用挂载外部或本地存储方式代替。
	DockerBaseSize *int32 `json:"dockerBaseSize,omitempty"`

	Runtime *Runtime `json:"runtime,omitempty"`
}

func (o ReinstallRuntimeConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReinstallRuntimeConfig struct{}"
	}

	return strings.Join([]string{"ReinstallRuntimeConfig", string(data)}, " ")
}
