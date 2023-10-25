package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTrafficMirrorFilterOption
type CreateTrafficMirrorFilterOption struct {

	// 功能说明：流量镜像筛选条件的描述信息 取值范围：0-255个字符，不能包含“<”和“>”
	Description *string `json:"description,omitempty"`

	// 功能说明：流量镜像筛选条件的名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`
}

func (o CreateTrafficMirrorFilterOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTrafficMirrorFilterOption struct{}"
	}

	return strings.Join([]string{"CreateTrafficMirrorFilterOption", string(data)}, " ")
}
