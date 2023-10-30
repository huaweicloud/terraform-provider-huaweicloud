package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Tags 标签列表。
type Tags struct {

	// 键。最大长度128个unicode字符。key不能为空。不能包含非打印字符ASCII(0-31)，“=”,“*”,“<”,“>”,“\\”,“,”,“|”,“/”
	Key *string `json:"key,omitempty"`

	// 值。每个值最大长度255个unicode字符，删除时如果value有值按照key/value删除，如果value没值，则按照key删除，可以为空字符串。,不能包含非打印字符ASCII(0-31), “=”,“*”,“<”,“>”,“\\”,“,”,“|”,“/”
	Value *string `json:"value,omitempty"`
}

func (o Tags) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Tags struct{}"
	}

	return strings.Join([]string{"Tags", string(data)}, " ")
}
