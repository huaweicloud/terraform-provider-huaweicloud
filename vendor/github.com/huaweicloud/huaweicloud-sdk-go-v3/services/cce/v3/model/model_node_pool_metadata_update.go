package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodePoolMetadataUpdate
type NodePoolMetadataUpdate struct {

	// 节点池名称。  > 命名规则： > >  - 以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围1-50位，且不能以中划线(-)结尾。 > >  - 不允许创建名为 DefaultPool 的节点池。
	Name string `json:"name"`
}

func (o NodePoolMetadataUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePoolMetadataUpdate struct{}"
	}

	return strings.Join([]string{"NodePoolMetadataUpdate", string(data)}, " ")
}
