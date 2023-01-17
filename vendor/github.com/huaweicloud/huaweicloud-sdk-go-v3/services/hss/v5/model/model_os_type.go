package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 操作系统类型，包含如下2种。   - Linux ：Linux。   - Windows ：Windows。
type OsType struct {
}

func (o OsType) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OsType struct{}"
	}

	return strings.Join([]string{"OsType", string(data)}, " ")
}
