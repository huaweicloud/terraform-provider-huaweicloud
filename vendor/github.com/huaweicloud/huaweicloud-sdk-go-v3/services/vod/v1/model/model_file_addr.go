package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 媒资存储参数信息。
type FileAddr struct {

	// OBS的bucket名称。
	Bucket string `json:"bucket"`

	// 桶所在的区域名， 如“华北-北京四”的区域名为“cn-north-4”，创建的桶所在区域必须和点播服务所在区域保持一致。
	Location string `json:"location"`

	// 文件的存储路径。
	Object string `json:"object"`
}

func (o FileAddr) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileAddr struct{}"
	}

	return strings.Join([]string{"FileAddr", string(data)}, " ")
}
