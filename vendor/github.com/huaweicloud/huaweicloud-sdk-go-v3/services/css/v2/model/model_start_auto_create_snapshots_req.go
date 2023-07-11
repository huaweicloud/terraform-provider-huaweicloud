package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartAutoCreateSnapshotsReq struct {

	// 指定要恢复的索引名称，多个索引用逗号隔开，默认恢复所有索引。支持使用“\\*”匹配多个索引，例如：2018-06\\*，表示恢复名称前缀是2018-06的所有索引的数据。 0～1024个字符，不能包含空格和大写字母，且不能包含\\\"\\\\<|>/?特殊字符。 默认值为\\*，表示恢复所有索引。
	Indices *string `json:"indices,omitempty"`

	// 设置快照保留的天数，范围是1～90。系统在半点时刻会自动删除超过保留天数的快照。
	Keepday int32 `json:"keepday"`

	// 每天创建快照的时刻，只支持整点，后面需加上时区，格式为“HH:mm z”，“HH:mm”表示整点时间，“z”表示时区。比如“00:00 GMT+08:00”、“01:00 GMT+08:00”等。
	Period string `json:"period"`

	// 自动创建的快照名称前缀，需要用户自己手动输入。只能包含1~32位小写字母、数字、中划线或者下划线，并且以小写字母开头。
	Prefix string `json:"prefix"`
}

func (o StartAutoCreateSnapshotsReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartAutoCreateSnapshotsReq struct{}"
	}

	return strings.Join([]string{"StartAutoCreateSnapshotsReq", string(data)}, " ")
}
