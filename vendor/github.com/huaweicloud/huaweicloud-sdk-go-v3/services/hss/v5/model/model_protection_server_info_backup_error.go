package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 备份错误信息
type ProtectionServerInfoBackupError struct {

	// 错误编码，包含如下2种。   - 0 ：无错误信息。   - 1 ：已綁定至其它存储库，无法开启备份。   - 2 ：备份库已超过最大限额。   - 3 ：CBR接口调用异常。
	ErrorCode *int32 `json:"error_code,omitempty"`

	// 错误描述
	ErrorDescription *string `json:"error_description,omitempty"`
}

func (o ProtectionServerInfoBackupError) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProtectionServerInfoBackupError struct{}"
	}

	return strings.Join([]string{"ProtectionServerInfoBackupError", string(data)}, " ")
}
