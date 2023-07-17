package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowResourceTagRequest Request Object
type ShowResourceTagRequest struct {

	// 资源ID
	ResourceId string `json:"resource_id"`

	// 项目ID，region级资源必选。
	ProjectId *string `json:"project_id,omitempty"`

	// 资源类型
	ResourceType string `json:"resource_type"`
}

func (o ShowResourceTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowResourceTagRequest struct{}"
	}

	return strings.Join([]string{"ShowResourceTagRequest", string(data)}, " ")
}
