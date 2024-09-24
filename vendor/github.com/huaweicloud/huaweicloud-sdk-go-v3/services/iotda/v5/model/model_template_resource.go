package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// TemplateResource 预调配模板设备资源结构体。
type TemplateResource struct {
	Device *DeviceResource `json:"device"`

	Policy *PolicyResource `json:"policy,omitempty"`
}

func (o TemplateResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateResource struct{}"
	}

	return strings.Join([]string{"TemplateResource", string(data)}, " ")
}
