package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteResourceInstanceTagRequest Request Object
type DeleteResourceInstanceTagRequest struct {

	// 由标签管理服务定义的资源类别，企业主机安全服务调用此接口时资源类别为hss
	ResourceType string `json:"resource_type"`

	// 由标签管理服务定义的资源id，企业主机安全服务调用此接口时资源id为配额ID
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
