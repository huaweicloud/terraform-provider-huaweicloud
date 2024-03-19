package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBinlogClearPolicyResponse Response Object
type ShowBinlogClearPolicyResponse struct {

	// binlog保留时长
	BinlogRetentionHours *int32 `json:"binlog_retention_hours,omitempty"`

	// 二进制日志保留策略,取值：time、fast - time:表示按时长保留二进制文件 - fast:表示快速清理,不保留二进制文件
	BinlogClearType *string `json:"binlog_clear_type,omitempty"`
	HttpStatusCode  int     `json:"-"`
}

func (o ShowBinlogClearPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBinlogClearPolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowBinlogClearPolicyResponse", string(data)}, " ")
}
