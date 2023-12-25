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
