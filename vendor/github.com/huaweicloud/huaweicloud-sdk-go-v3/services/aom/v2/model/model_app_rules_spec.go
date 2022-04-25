package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则详情。
type AppRulesSpec struct {

	// 服务类型,用于标记服务的分类,仅用于规则分类和界面展示。可以填写任意字段,如按技术栈分类可填写Java,Python。按作用分类可填写collector(采集),database(数据库)等。
	AppType string `json:"appType"`

	// 属性列表(暂不使用,可不传)。 cmdLine、env
	AttrList *[]string `json:"attrList,omitempty"`

	// 是否开启日志采集。 true、false
	DetectLog string `json:"detectLog"`

	// 规则发现部分,数组中有多个对象时表示需要同时满足所有条件的进程才会被匹配到。 checkType为cmdLine时checkMode填contain,checkContent格式为[“xxx”]表示进程命令行参数中需要包含xxx。checkType为env时checkMode填contain,checkContent格式为 [\"k1\",\"v1\"]表示进程环境变量中需要包含名为k1值为v1的环境变量。checkType为scope时checkMode填equals,checkContent格式为节点ID数组[\"hostId1”,”hostId2”],表示规则仅会在这些节点上生效(如果不指定节点范围,规则将下发到该项目所有的节点)。
	DiscoveryRule []DiscoveryRule `json:"discoveryRule"`

	// 是否为默认规则。 true、false
	IsDefaultRule string `json:"isDefaultRule"`

	// 是否为规则预探测场景(预探测场景不会保存规则,仅用于规则下发之前验证该规则能否有效发现节点上的进程)。 true、false
	IsDetect string `json:"isDetect"`

	// 日志文件的后缀。 log、trace、out
	LogFileFix []string `json:"logFileFix"`

	// 日志路径配置规则。 当cmdLineHash为固定字符串时,指定日志路径或者日志文件。否则只采集进程当前打开的以.log和.trace结尾的文件。nameType取值cmdLineHash时,args格式为[\"00001\"],value格式为[\"/xxx/xx.log\"],表示当启动命令是00001时,日志路径为/xxx/xx.log。
	LogPathRule *[]LogPathRule `json:"logPathRule,omitempty"`

	NameRule *NameRule `json:"nameRule"`

	// 规则优先级。 1~9999的整数字符串,默认取值为9999
	Priority int32 `json:"priority"`
}

func (o AppRulesSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppRulesSpec struct{}"
	}

	return strings.Join([]string{"AppRulesSpec", string(data)}, " ")
}
