package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// GetTargetImageIdDetail 镜像详情信息。
type GetTargetImageIdDetail struct {

	// 可以升级的目标镜像ID。
	Id *string `json:"id,omitempty"`

	// 可以升级的目标镜像名称。
	DisplayName *string `json:"displayName,omitempty"`

	// 镜像描述信息。
	ImageDesc *string `json:"imageDesc,omitempty"`

	// 镜像引擎类型。
	DatastoreType *string `json:"datastoreType,omitempty"`

	// 镜像引擎版本。
	DatastoreVersion *string `json:"datastoreVersion,omitempty"`

	// 优先级。
	Priority *int32 `json:"priority,omitempty"`
}

func (o GetTargetImageIdDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetTargetImageIdDetail struct{}"
	}

	return strings.Join([]string{"GetTargetImageIdDetail", string(data)}, " ")
}
