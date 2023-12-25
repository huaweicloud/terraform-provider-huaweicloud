package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeIsolatedFileRequestInfo 恢复已隔离的文件详情
type ChangeIsolatedFileRequestInfo struct {

	// 需要恢复的文件列表
	DataList *[]IsolatedFileRequestInfo `json:"data_list,omitempty"`
}

func (o ChangeIsolatedFileRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeIsolatedFileRequestInfo struct{}"
	}

	return strings.Join([]string{"ChangeIsolatedFileRequestInfo", string(data)}, " ")
}
