package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListShareBackupsRequest Request Object
type ListShareBackupsRequest struct {

	// 实例ID。
	InstanceId *string `json:"instance_id,omitempty"`

	// 实例名称。
	InstanceName *string `json:"instance_name,omitempty"`

	// 备份名称。
	BackupName *string `json:"backup_name,omitempty"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *string `json:"offset,omitempty"`

	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100。
	Limit *string `json:"limit,omitempty"`
}

func (o ListShareBackupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListShareBackupsRequest struct{}"
	}

	return strings.Join([]string{"ListShareBackupsRequest", string(data)}, " ")
}
