package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlBackupListRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID。

	InstanceId *string `json:"instance_id,omitempty"`
	// 备份ID。

	BackupId *string `json:"backup_id,omitempty"`
	// 备份类型，取值：   \"auto\"：自动全量备份   \"manual\"：手动全量备份

	BackupType *string `json:"backup_type,omitempty"`
	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。

	Offset *string `json:"offset,omitempty"`
	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100。

	Limit *string `json:"limit,omitempty"`
	// 查询开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	BeginTime *string `json:"begin_time,omitempty"`
	// 查询结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”，且大于查询开始时间。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。

	EndTime *string `json:"end_time,omitempty"`
}

func (o ShowGaussMySqlBackupListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlBackupListRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlBackupListRequest", string(data)}, " ")
}
