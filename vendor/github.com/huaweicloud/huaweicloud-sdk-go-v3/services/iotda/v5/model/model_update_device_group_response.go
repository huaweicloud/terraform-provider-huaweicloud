package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateDeviceGroupResponse struct {

	// 设备组ID，用于唯一标识一个设备组，在创建设备组时由物联网平台分配。
	GroupId *string `json:"group_id,omitempty"`

	// 设备组名称，单个资源空间下不可重复。
	Name *string `json:"name,omitempty"`

	// 设备组描述。
	Description *string `json:"description,omitempty"`

	// 父设备组ID，该设备组的父设备组ID。
	SuperGroupId   *string `json:"super_group_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDeviceGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDeviceGroupResponse struct{}"
	}

	return strings.Join([]string{"UpdateDeviceGroupResponse", string(data)}, " ")
}
