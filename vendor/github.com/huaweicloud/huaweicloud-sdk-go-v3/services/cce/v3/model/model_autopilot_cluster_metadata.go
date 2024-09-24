package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AutopilotClusterMetadata 支持通过 annotations[\"cluster.install.addons/install\"] 来指定创建集群时需要安装的插件
type AutopilotClusterMetadata struct {

	// 集群名称。  命名规则：以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围4-128位，且不能以中划线(-)结尾。
	Name string `json:"name"`

	// 集群ID，资源唯一标识，创建成功后自动生成，填写无效。在创建包周期集群时，响应体不返回集群ID。
	Uid *string `json:"uid,omitempty"`

	// 集群显示名，用于在 CCE 界面显示，该名称创建后可修改。  命名规则：以小写字母开头，由小写字母、数字、中划线(-)组成，长度范围4-128位，且不能以中划线(-)结尾。  显示名和其他集群的名称、显示名不可以重复。  在创建集群、更新集群请求体中，集群显示名alias未指定或取值为空，表示与集群名称name一致。在查询集群等响应体中，集群显示名alias将必然返回，未配置时将返回集群名称name。
	Alias *string `json:"alias,omitempty"`

	// 集群注解，由key/value组成：  ``` \"annotations\": {    \"key1\" : \"value1\",    \"key2\" : \"value2\" } ```  >    - Annotations不用于标识和选择对象。Annotations中的元数据可以是small或large，structured或unstructured，并且可以包括标签不允许使用的字符。 >    - 该字段不会被数据库保存，当前仅用于指定集群待安装插件。
	Annotations map[string]string `json:"annotations,omitempty"`

	// 集群标签，key/value对格式。  >  该字段值由系统自动生成，用于升级时前端识别集群支持的特性开关，用户指定无效。
	Labels map[string]string `json:"labels,omitempty"`

	// 集群创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 集群更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`
}

func (o AutopilotClusterMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AutopilotClusterMetadata struct{}"
	}

	return strings.Join([]string{"AutopilotClusterMetadata", string(data)}, " ")
}
