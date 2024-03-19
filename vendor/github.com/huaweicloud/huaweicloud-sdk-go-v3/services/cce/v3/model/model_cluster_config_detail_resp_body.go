package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterConfigDetailRespBody 获取指定集群配置项列表返回体
type ClusterConfigDetailRespBody struct {

	// 配置参数，由key/value组成。
	KubeApiserver *[]PackageOptions `json:"kube-apiserver,omitempty"`
}

func (o ClusterConfigDetailRespBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterConfigDetailRespBody struct{}"
	}

	return strings.Join([]string{"ClusterConfigDetailRespBody", string(data)}, " ")
}
