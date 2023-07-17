package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestoreSnapshotReq struct {

	// 快照要恢复到的集群的ID。
	TargetCluster string `json:"targetCluster"`

	// 指定要恢复的索引名称，多个索引用逗号隔开，默认恢复所有索引。支持使用“\\*”匹配多个索引，例如：2018-06\\*，表示恢复名称前缀是2018-06的所有索引的数据。  0～1024个字符，不能包含空格和大写字母，且不能包含\\\"\\\\<|>/?特殊字符。
	Indices *string `json:"indices,omitempty"`

	// 匹配要恢复的索引规则，最大支持1024个字符。根据此处定义的过滤条件去恢复符合条件的索引，过滤条件请使用正则表达式。  0～1024个字符，不能包含空格和大写字母，且不能包含\\\"\\\\<|>/?,特殊字符。   renamePattern参数与renameReplacement参数必须同时设置才能生效。
	RenamePattern *string `json:"renamePattern,omitempty"`

	// 索引重命名的规则。0～1024个字符，不能包含空格和大写字母，且不能包含\\\"\\\\<|>/?,特殊字符。例如，“restored_index_$1”表示在所有恢复的索引名称前面加上“restored_”。    renamePattern参数与renameReplacement参数必须同时设置才能生效。
	RenameReplacement *string `json:"renameReplacement,omitempty"`
}

func (o RestoreSnapshotReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreSnapshotReq struct{}"
	}

	return strings.Join([]string{"RestoreSnapshotReq", string(data)}, " ")
}
