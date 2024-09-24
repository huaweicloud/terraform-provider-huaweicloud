package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotClusterInformationSpec
type AutopilotClusterInformationSpec struct {

	// 集群的描述信息。  1. 字符取值范围[0,200]。不包含~$%^&*<>[]{}()'\"#\\等特殊字符。 2. 仅运行（Available）的集群允许修改。
	Description *string `json:"description,omitempty"`

	// 集群的API Server服务端证书中的自定义SAN（Subject Alternative Name）字段，遵从SSL标准X509定义的格式规范。  1. 不允许出现同名重复。 2. 格式符合IP和域名格式。  示例: ``` SAN 1: DNS Name=example.com SAN 2: DNS Name=www.example.com SAN 3: DNS Name=example.net SAN 4: IP Address=93.184.216.34 ```
	CustomSan *[]string `json:"customSan,omitempty"`

	EniNetwork *AutopilotEniNetworkUpdate `json:"eniNetwork,omitempty"`
}

func (o AutopilotClusterInformationSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterInformationSpec struct{}"
	}

	return strings.Join([]string{"AutopilotClusterInformationSpec", string(data)}, " ")
}
