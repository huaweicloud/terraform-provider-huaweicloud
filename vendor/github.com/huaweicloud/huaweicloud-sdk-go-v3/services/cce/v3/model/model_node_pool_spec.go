package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// NodePoolSpec
type NodePoolSpec struct {

	// 节点池类型。不填写时默认为vm。  - vm：弹性云服务器 - ElasticBMS：C6型弹性裸金属通用计算增强型云服务器，规格示例：c6.22xlarge.2.physical - pm: 裸金属服务器
	Type *NodePoolSpecType `json:"type,omitempty"`

	NodeTemplate *NodeSpec `json:"nodeTemplate"`

	// 节点池初始化节点个数。查询时为节点池目标节点数量。
	InitialNodeCount *int32 `json:"initialNodeCount,omitempty"`

	Autoscaling *NodePoolNodeAutoscaling `json:"autoscaling,omitempty"`

	NodeManagement *NodeManagement `json:"nodeManagement,omitempty"`

	// 1.21版本集群节点池支持绑定安全组，最多五个。
	PodSecurityGroups *[]SecurityId `json:"podSecurityGroups,omitempty"`

	// 节点池扩展伸缩组配置列表，详情参见ExtensionScaleGroup类型定义
	ExtensionScaleGroups *[]ExtensionScaleGroup `json:"extensionScaleGroups,omitempty"`

	// 节点池自定义安全组相关配置。支持节点池新扩容节点绑定指定的安全组。  - 未指定安全组ID，新建节点将添加Node节点默认安全组。  - 指定有效安全组ID，新建节点将使用指定安全组。  - 指定安全组，应避免对CCE运行依赖的端口规则进行修改。[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/cce_faq/cce_faq_00265.html)。](tag:hws)[详细设置请参考[集群安全组规则配置](https://support.huaweicloud.com/intl/zh-cn/cce_faq/cce_faq_00265.html)。](tag:hws_hk)
	CustomSecurityGroups *[]string `json:"customSecurityGroups,omitempty"`
}

func (o NodePoolSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodePoolSpec struct{}"
	}

	return strings.Join([]string{"NodePoolSpec", string(data)}, " ")
}

type NodePoolSpecType struct {
	value string
}

type NodePoolSpecTypeEnum struct {
	VM          NodePoolSpecType
	ELASTIC_BMS NodePoolSpecType
	PM          NodePoolSpecType
}

func GetNodePoolSpecTypeEnum() NodePoolSpecTypeEnum {
	return NodePoolSpecTypeEnum{
		VM: NodePoolSpecType{
			value: "vm",
		},
		ELASTIC_BMS: NodePoolSpecType{
			value: "ElasticBMS",
		},
		PM: NodePoolSpecType{
			value: "pm",
		},
	}
}

func (c NodePoolSpecType) Value() string {
	return c.value
}

func (c NodePoolSpecType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *NodePoolSpecType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
