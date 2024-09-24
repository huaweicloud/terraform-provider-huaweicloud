package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterListDatastore 数据搜索引擎类型。
type ClusterListDatastore struct {

	// 引擎类型，目前只支持elasticsearch。
	Type *string `json:"type,omitempty"`

	// CSS集群引擎版本号。详细请参考CSS[支持的集群版本](css_03_0056.xml)。
	Version *string `json:"version,omitempty"`

	// 是否支持安全模式
	SupportSecuritymode *bool `json:"supportSecuritymode,omitempty"`
}

func (o ClusterListDatastore) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterListDatastore struct{}"
	}

	return strings.Join([]string{"ClusterListDatastore", string(data)}, " ")
}
