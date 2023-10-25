package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TrafficMirrorFilter
type TrafficMirrorFilter struct {

	// 功能说明：流量镜像筛选条件ID
	Id string `json:"id"`

	// 功能说明：项目ID
	ProjectId string `json:"project_id"`

	// 功能说明：流量镜像筛选条件的描述信息 取值范围：0-255个字符，不能包含“<”和“>”
	Description string `json:"description"`

	// 功能说明：流量镜像筛选条件的名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：入方向筛选规则列表
	IngressRules []TrafficMirrorFilterRule `json:"ingress_rules"`

	// 功能说明：出方向筛选规则列表
	EgressRules []TrafficMirrorFilterRule `json:"egress_rules"`

	// 创建时间戳
	CreatedAt *sdktime.SdkTime `json:"created_at"`

	// 更新时间戳
	UpdatedAt *sdktime.SdkTime `json:"updated_at"`
}

func (o TrafficMirrorFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TrafficMirrorFilter struct{}"
	}

	return strings.Join([]string{"TrafficMirrorFilter", string(data)}, " ")
}
