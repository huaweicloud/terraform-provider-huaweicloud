package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TraceResource struct {

	// 云服务类型。必须为已对接CTS的云服务的英文缩写，且服务类型一般为大写字母。
	ServiceType *string `json:"service_type,omitempty"`

	// 云服务对应的资源类型列表。
	Resource *[]string `json:"resource,omitempty"`
}

func (o TraceResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TraceResource struct{}"
	}

	return strings.Join([]string{"TraceResource", string(data)}, " ")
}
