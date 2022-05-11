package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 应用命名部分。 nameType取值cmdLine时args格式为[\"start\",\"end\"],表示抽取命令行中start、end之间的字符。 nameType取值cmdLine时args格式为 [\"aa\"],表示抽取环境变量名为aa对应的环境变量值。 nameType取值str时,args格式为[\"fix\"],表示服务名称最后拼接固定文字fix。 nameType取值cmdLineHash时,args格式为[\"0001\"],value格式为[\"ser\"],表示当启动命令是0001时,应用名称为ser。
type ApplicationNameRule struct {

	// 取值类型。 从cmdLineHash、cmdLine、env、str里面选取。
	NameType string `json:"nameType"`

	// 输入值。
	Args []string `json:"args"`

	// 服务名(仅nameType为cmdLineHash时填写)。
	Value *[]string `json:"value,omitempty"`
}

func (o ApplicationNameRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ApplicationNameRule struct{}"
	}

	return strings.Join([]string{"ApplicationNameRule", string(data)}, " ")
}
