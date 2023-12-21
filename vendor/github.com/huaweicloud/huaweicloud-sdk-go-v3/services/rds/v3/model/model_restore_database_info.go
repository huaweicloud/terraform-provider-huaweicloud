package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreDatabaseInfo 库级恢复数据库信息
type RestoreDatabaseInfo struct {

	// 恢复前库名
	OldName string `json:"old_name"`

	// 恢复后库名
	NewName string `json:"new_name"`
}

func (o RestoreDatabaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreDatabaseInfo struct{}"
	}

	return strings.Join([]string{"RestoreDatabaseInfo", string(data)}, " ")
}
