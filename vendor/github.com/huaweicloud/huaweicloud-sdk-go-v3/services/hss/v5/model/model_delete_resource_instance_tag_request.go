package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteResourceInstanceTagRequest Request Object
type DeleteResourceInstanceTagRequest struct {

	// 资源类别，hss
	ResourceType string `json:"resource_type"`

	// 资源ID
	ResourceId string `json:"resource_id"`

	// 待删除的key
	Key string `json:"key"`
}

func (o DeleteResourceInstanceTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteResourceInstanceTagRequest struct{}"
	}

	return strings.Join([]string{"DeleteResourceInstanceTagRequest", string(data)}, " ")
}
