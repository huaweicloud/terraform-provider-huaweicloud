package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateSnapshotReq struct {

	// 快照名称，快照名称在4位到64位之间，必须以字母开头，可以包含字母、数字、中划线或者下划线，注意字母不能大写且不能包含其他特殊字符。
	Name string `json:"name"`

	// 快照描述，0～256个字符，不能包含!<>=&\\\"'字符。
	Description *string `json:"description,omitempty"`

	// 指定要备份的索引名称，多个索引用逗号隔开，默认备份所有索引。支持使用“\\*”匹配多个索引，例如：2018-06\\*，表示备份名称前缀是2018-06的所有索引的数据。  0～1024个字符，不能包含空格和大写字母，且不能包含\\\"\\\\<|>/?特殊字符。
	Indices *string `json:"indices,omitempty"`
}

func (o CreateSnapshotReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSnapshotReq struct{}"
	}

	return strings.Join([]string{"CreateSnapshotReq", string(data)}, " ")
}
