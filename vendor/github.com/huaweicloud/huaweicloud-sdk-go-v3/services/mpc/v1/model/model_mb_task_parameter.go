package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MbTaskParameter struct {

	// 具体状态描述，FAILED时可用于分析问题。
	StatusDescription *string `json:"status_description,omitempty"`

	// 输出文件名称。
	OutputFilename *string `json:"output_filename,omitempty"`

	Metadata *MetaData `json:"metadata,omitempty"`
}

func (o MbTaskParameter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MbTaskParameter struct{}"
	}

	return strings.Join([]string{"MbTaskParameter", string(data)}, " ")
}
