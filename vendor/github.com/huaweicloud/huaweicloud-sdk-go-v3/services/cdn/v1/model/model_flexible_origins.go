package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FlexibleOrigins 灵活回源信息,最多20条。
type FlexibleOrigins struct {

	// URI的匹配方式，支持文件后缀（file_extension）和路径前缀（file_path）。
	MatchType string `json:"match_type"`

	// file_extension（文件后缀）： 支持所有格式的文件类型。 输入首字符为“.”，以“;”进行分隔。 输入的文件后缀名总数不能超过20个。 file_path（目录路径）：输入要求以“/”作为首字符，以“;”进行分隔，输入的目录路径总数不能超过20个。
	MatchPattern string `json:"match_pattern"`

	// 优先级取值范围为1~100，数值越大优先级越高。
	Priority int32 `json:"priority"`

	// 回源信息。  > 每个目录的回源源站数量不超过1个。
	BackSources []BackSources `json:"back_sources"`
}

func (o FlexibleOrigins) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FlexibleOrigins struct{}"
	}

	return strings.Join([]string{"FlexibleOrigins", string(data)}, " ")
}
