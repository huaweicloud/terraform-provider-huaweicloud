package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Metadata 资源描述基本信息，集合类的元素类型，包含一组由不同名称定义的属性。
type Metadata struct {

	// 唯一id标识
	Uid *string `json:"uid,omitempty"`

	// 资源名称
	Name *string `json:"name,omitempty"`

	// 资源标签，key/value对格式，接口保留字段，填写不会生效
	Labels map[string]string `json:"labels,omitempty"`

	// 资源注解，由key/value组成
	Annotations map[string]string `json:"annotations,omitempty"`

	// 更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`

	// 创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`
}

func (o Metadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Metadata struct{}"
	}

	return strings.Join([]string{"Metadata", string(data)}, " ")
}
