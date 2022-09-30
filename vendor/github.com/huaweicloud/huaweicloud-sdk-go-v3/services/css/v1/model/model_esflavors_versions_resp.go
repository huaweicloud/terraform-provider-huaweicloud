package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 引擎列表。
type EsflavorsVersionsResp struct {

	// Esasticsearch引擎版本号。详细请参考CSS[支持的集群版本](css_03_0056.xml)。
	Version *string `json:"version,omitempty"`

	// 规格列表。
	Flavors *[]EsflavorsVersionsFlavorsResp `json:"flavors,omitempty"`

	// 实例类型，包括为ess、ess-cold、ess-master和ess-client。
	Type *string `json:"type,omitempty"`
}

func (o EsflavorsVersionsResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EsflavorsVersionsResp struct{}"
	}

	return strings.Join([]string{"EsflavorsVersionsResp", string(data)}, " ")
}
