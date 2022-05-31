package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsObjInfo struct {

	// OBS的bucket名称。
	Bucket string `json:"bucket"`

	// OBS桶所在的区域，且必须与使用的MPC区域保持一致。
	Location string `json:"location"`

	// OBS对象路径，遵守OSS Object定义。  - 当用于指示input时,需要指定到具体对象。 - 当用于指示output时, 只需指定到转码结果期望存放的路径。
	Object string `json:"object"`

	// 文件名，仅用于转封装指定输出名称。  - 当指定了此参数时，输出的对象名为object/file_name 。 - 当不指定此参数时，输出的对象名为object/xxx，其中xxx由MPC指定。
	FileName *string `json:"file_name,omitempty"`
}

func (o ObsObjInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsObjInfo struct{}"
	}

	return strings.Join([]string{"ObsObjInfo", string(data)}, " ")
}
