package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreDatabaseRequest 库级恢复请求信息
type PostgreSqlRestoreDatabaseRequest struct {

	// 库级恢复实例信息
	Instances []RestoreDatabaseInstance `json:"instances"`
}

func (o PostgreSqlRestoreDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreDatabaseRequest struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreDatabaseRequest", string(data)}, " ")
}
