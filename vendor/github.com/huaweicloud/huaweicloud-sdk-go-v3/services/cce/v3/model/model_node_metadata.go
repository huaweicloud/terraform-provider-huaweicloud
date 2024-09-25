package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeMetadata
type NodeMetadata struct {

	// 节点名称 > 命名规则：以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围1-56位，且不能以中划线(-)结尾。 > 若name未指定或指定为空字符串，则按照默认规则生成节点名称。默认规则为：“集群名称-随机字符串”，若集群名称过长，则只取前36个字符。 > 若节点数量(count)大于1时，则按照默认规则会在用户输入的节点名称末尾添加随机字符串。默认规则为：“用户输入名称-随机字符串”，若用户输入的节点名称长度范围超过50位时，系统截取前50位，并在末尾添加随机字符串。
	Name *string `json:"name,omitempty"`

	// 节点ID，资源唯一标识，创建成功后自动生成，填写无效
	Uid *string `json:"uid,omitempty"`

	// CCE自有节点标签，非Kubernetes原生labels。  标签可用于选择对象并查找满足某些条件的对象集合，格式为key/value键值对。  示例：  ``` \"labels\": {   \"key\" : \"value\" } ```
	Labels map[string]string `json:"labels,omitempty"`

	// CCE自有节点注解，非Kubernetes原生annotations，格式为key/value键值对。 示例： ``` \"annotations\": {   \"key1\" : \"value1\",   \"key2\" : \"value2\" } ``` > - Annotations不用于标识和选择对象。Annotations中的元数据可以是small或large，structured或unstructured，并且可以包括标签不允许使用的字符。 > - 仅用于查询，不支持请求时传入，填写无效。
	Annotations map[string]string `json:"annotations,omitempty"`

	// 创建时间，创建成功后自动生成，填写无效
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 更新时间，创建成功后自动生成，填写无效
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`

	OwnerReference *NodeOwnerReference `json:"ownerReference,omitempty"`
}

func (o NodeMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeMetadata struct{}"
	}

	return strings.Join([]string{"NodeMetadata", string(data)}, " ")
}
