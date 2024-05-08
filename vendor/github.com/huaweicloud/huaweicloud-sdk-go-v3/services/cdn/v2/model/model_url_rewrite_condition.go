package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UrlRewriteCondition 匹配条件。
type UrlRewriteCondition struct {

	// 匹配类型。   - catalog：指定目录下的文件需执行访问URL重写规则，   - full_path：某个完整路径下的文件需执行访问URL重写规则。
	MatchType string `json:"match_type"`

	// 匹配内容。当match_type为catalog时，为目录路径，输入要求以“/”作为首字符，以“,”进行分隔，如/test/folder01,/test/folder02，并且输入的目录路径总数不超过20个。 当match_type为full_path时，为全路径，输入要求以“/”作为首字符，支持匹配指定目录下的具体文件，或者带通配符“\\*”的文件，单条全路径缓存规则里仅支持配置一个全路径，如/test/index.html或/test/\\*.jpg。
	MatchValue string `json:"match_value"`

	// 访问URL重写规则的优先级。取值为1-100之间的整数，数值越大优先级越高。优先级设置具有唯一性，不支持多条规则设置同一优先级，且优先级不能为空。
	Priority int32 `json:"priority"`
}

func (o UrlRewriteCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UrlRewriteCondition struct{}"
	}

	return strings.Join([]string{"UrlRewriteCondition", string(data)}, " ")
}
