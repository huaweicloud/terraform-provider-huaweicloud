package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IsolatedFileResponseInfo 已隔离文件详情
type IsolatedFileResponseInfo struct {

	// 服务器ID
	HostId string `json:"host_id"`

	// 服务器名称
	HostName string `json:"host_name"`

	// 文件哈希
	FileHash string `json:"file_hash"`

	// 文件路径
	FilePath string `json:"file_path"`

	// 隔离状态，包含如下:   - isolated : 已隔离   - restored : 已恢复   - isolating : 已下发隔离任务   - restoring : 已下发恢复任务
	IsolationStatus string `json:"isolation_status"`

	// 文件属性
	FileAttr string `json:"file_attr"`

	// 更新时间，毫秒
	UpdateTime int64 `json:"update_time"`
}

func (o IsolatedFileResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IsolatedFileResponseInfo struct{}"
	}

	return strings.Join([]string{"IsolatedFileResponseInfo", string(data)}, " ")
}
