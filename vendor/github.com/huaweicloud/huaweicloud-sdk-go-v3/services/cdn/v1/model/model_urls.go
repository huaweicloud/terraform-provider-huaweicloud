package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 具体url信息
type Urls struct {

	// urlid
	Id *int32 `json:"id,omitempty"`

	// url具体值
	Url *string `json:"url,omitempty"`

	// url状态
	Status *string `json:"status,omitempty"`

	// 任务类型
	Type *string `json:"type,omitempty"`

	// 任务id
	TaskId *int64 `json:"task_id,omitempty"`

	// 修改时间戳（毫秒）
	ModifyTime *int32 `json:"modify_time,omitempty"`

	// 创建时间戳（毫秒）
	CreateTime *int32 `json:"create_time,omitempty"`

	// 文件类型，目录还是文件
	FileType *string `json:"file_type,omitempty"`
}

func (o Urls) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Urls struct{}"
	}

	return strings.Join([]string{"Urls", string(data)}, " ")
}
