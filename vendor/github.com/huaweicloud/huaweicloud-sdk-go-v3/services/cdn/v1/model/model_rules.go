package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Rules struct {

	// 0：全部类型，表示匹配所有文件，默认值。  1：文件类型，表示按文件后缀匹配。  2：文件夹类型，表示按目录匹配。  3：文件全路径类型，表示按文件全路径匹配。
	RuleType int32 `json:"rule_type"`

	// 缓存匹配设置。  当rule_type为0时，为空。  当rule_type为1时，为文件后缀，输入首字符为“.”，以“;”进行分隔，如.jpg;.zip;.exe，并且输入的文件名后缀总数不超过20个。 当rule_type为2时，为目录，输入要求以“/”作为首字符，以“;”进行分隔，如/test/folder01;/test/folder02，并且输入的目录路径总数不超过20个。 当rule_type为3时，为全路径，输入要求以“/”作为首字符，支持匹配指定目录下的具体文件，或者带通配符“*”的文件，如/test/index.html或/test/_*.jpg。
	Content *string `json:"content,omitempty"`

	// 缓存时间。最大支持365天。
	Ttl int32 `json:"ttl"`

	// 缓存时间单位。1：秒；2：分；3：小时；4：天。
	TtlType int32 `json:"ttl_type"`

	// 此条配置的权重值, 默认值1，数值越大，优先级越高。取值范围为1-100，权重值不能相同。
	Priority int32 `json:"priority"`
}

func (o Rules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Rules struct{}"
	}

	return strings.Join([]string{"Rules", string(data)}, " ")
}
