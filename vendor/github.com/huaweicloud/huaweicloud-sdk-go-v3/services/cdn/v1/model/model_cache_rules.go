package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CacheRules struct {

	// 类型,all：匹配所有文件， file_extension：按文件后缀匹配， catalog：按目录匹配， full_path：全路径匹配， home_page：按首页匹配。
	MatchType string `json:"match_type"`

	// 缓存匹配设置， 当match_type为all时，为空。当match_type为file_extension时，为文件后缀，输入首字符为“.”，以“,”进行分隔， 如.jpg,.zip,.exe，并且输入的文 件名后缀总数不超过20个。 当match_type为catalog时，为目录，输入要求以“/”作为首字符， 以“,”进行分隔，如/test/folder01,/test/folder02，并且输入的目录路径总数不超过20个。  当match_type为full_path时，为全路径，输入要求以“/”作为首字符，支持匹配指定目录下的具体文件，或者带通配符“\\*”的文件， 如/test/index.html,/test/\\*.jpg。 当match_type为home_page时，为空。
	MatchValue *string `json:"match_value,omitempty"`

	// 缓存过期时间，最大支持365天。
	Ttl int32 `json:"ttl"`

	// 缓存过期时间单位，s：秒；m：分；h：小时；d：天。
	TtlUnit string `json:"ttl_unit"`

	// 此条配置的优先级, 默认值1，数值越大，优先级越高，取值范围为1-100，优先级不能相同。
	Priority int32 `json:"priority"`

	// 缓存遵循源站开关，on：打开，off：关闭。
	FollowOrigin string `json:"follow_origin"`

	// URL参数， del_params：忽略指定URL参数， reserve_params：保留指定URL参数， ignore_url_params：忽略全部URL参数， full_url：使用完整URL参数。
	UrlParameterType string `json:"url_parameter_type"`

	// URL参数值，最多设置10条，以\",\"分隔。
	UrlParameterValue *string `json:"url_parameter_value,omitempty"`
}

func (o CacheRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CacheRules struct{}"
	}

	return strings.Join([]string{"CacheRules", string(data)}, " ")
}
