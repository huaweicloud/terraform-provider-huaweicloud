package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 列表文件配置。
type ListFile struct {

	// 对象列表文件或URL列表文件对象名。
	ListFileKey string `json:"list_file_key"`

	// 存放对象列表文件的OBS桶名。  请确保与目的端桶处于同一区域，否则将导致任务创建失败。
	ObsBucket string `json:"obs_bucket"`
}

func (o ListFile) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListFile struct{}"
	}

	return strings.Join([]string{"ListFile", string(data)}, " ")
}
