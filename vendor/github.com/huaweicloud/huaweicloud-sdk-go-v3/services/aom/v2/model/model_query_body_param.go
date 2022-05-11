package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type QueryBodyParam struct {

	// 取值范围 app_log,node _log,custom_log 日志类型字段:app_log:应用日志 node_log:主机日志 custom_log:自定义配置路径日志。
	Category string `json:"category"`

	// 搜索结束时间(UTC时间，毫秒级)。
	EndTime int64 `json:"endTime"`

	// 取值范围 0、1 。搜索时是否隐藏系统日志，默认0为隐藏1为显示。
	HideSyslog *int32 `json:"hideSyslog,omitempty"`

	// 1.支持关键词精确搜索。关键词指相邻两个分词符之间的单词。 2.支持关键词模糊匹配搜索，例如输入“RROR”或“ERRO?”或“*ROR*”或“ERR*”或“ER*OR”。 3.支持短语精确搜索，例如输入“Start to refresh alm Statistic”。 4.支持关键词的“与”、“或”组合搜索。格式为“query&&logs”或“query||logs”。 说明： 当前默认分词符有  , '\";=()[]{}@&<>/:\\n\\t\\r
	KeyWord *string `json:"keyWord,omitempty"`

	// 日志单行序列号第一次查询时不需要此参数，后续分页查询时需要使用可从上次查询的返回信息中获取.
	LineNum *string `json:"lineNum,omitempty"`

	// 表示每次查询的日志条数不填时默认为5000，建议您设置为100。 第一次查询时使用pageSize 后续分页查询时使用size。
	PageSizeSize *string `json:"pageSize/size,omitempty"`

	SearchKey *SearchKey `json:"searchKey"`

	// 搜索起始时间(UTC时间，毫秒级)。
	StartTime int64 `json:"startTime"`

	// 表示此次查询为分页查询，第一次查询时不需要此参数，后续分页查询时需要使用。
	Type *string `json:"type,omitempty"`

	// 标识按照lineNum升序查询还是降序查询。  true：降序（lineNum由大到小：时间从新到老）。  false：升序（lineNum由小到大：即时间从老到新）。
	IsDesc *bool `json:"isDesc,omitempty"`
}

func (o QueryBodyParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryBodyParam struct{}"
	}

	return strings.Join([]string{"QueryBodyParam", string(data)}, " ")
}
