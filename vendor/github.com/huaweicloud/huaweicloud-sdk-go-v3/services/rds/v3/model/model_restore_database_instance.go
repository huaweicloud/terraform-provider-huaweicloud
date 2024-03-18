package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreDatabaseInstance 查询可恢复库的响应信息
type RestoreDatabaseInstance struct {

	// 恢复时间
	RestoreTime int64 `json:"restore_time"`

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 是否使用极速恢复，可先根据”获取实例是否能使用极速恢复“接口判断本次恢复是否能使用极速恢复。 如果实例使用了XA事务，采用极速恢复的方式会导致恢复失败！
	IsFastRestore *bool `json:"is_fast_restore,omitempty"`

	// 库信息
	Databases []RestoreDatabaseInfo `json:"databases"`
}

func (o RestoreDatabaseInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreDatabaseInstance struct{}"
	}

	return strings.Join([]string{"RestoreDatabaseInstance", string(data)}, " ")
}
