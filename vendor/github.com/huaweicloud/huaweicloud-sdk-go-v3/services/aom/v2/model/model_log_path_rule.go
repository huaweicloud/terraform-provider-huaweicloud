package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 日志路径配置规则。 当cmdLineHash为固定字符串时,指定日志路径或者日志文件。否则只采集进程当前打开的以.log和.trace结尾的文件。nameType取值cmdLineHash时,args格式为[\"00001\"],value格式为[\"/xxx/xx.log\"],表示当启动命令是00001时,日志路径为/xxx/xx.log。
type LogPathRule struct {

	// 命令行。
	Args []string `json:"args"`

	// 取值类型。 cmdLineHash
	NameType string `json:"nameType"`

	// 日志路径。
	Value []string `json:"value"`
}

func (o LogPathRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LogPathRule struct{}"
	}

	return strings.Join([]string{"LogPathRule", string(data)}, " ")
}
