package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 服务命名部分,数组中有多个对象时表示将每个对象抽取到的字符串拼接作为服务的名称。 nameType取值cmdLine时args格式为[\"start\",\"end\"],表示抽取命令行中start、end之间的字符。 nameType取值cmdLine时args格式为[\"aa\"],表示抽取环境变量名为aa对应的环境变量值。 nameType取值str时,args格式为[\"fix\"],表示服务名称最后拼接固定文字fix。 nameType取值cmdLineHash时,args格式为[\"0001\"],value格式为[\"ser\"],表示当启动命令是0001时,服务名称为ser。
type AppNameRule struct {

	// 取值类型。 从cmdLineHash、cmdLine、env、str里面选取。
	NameType string `json:"nameType"`

	// 输入值。
	Args []string `json:"args"`

	// 服务名(仅nameType为cmdLineHash时填写)。
	Value *[]string `json:"value,omitempty"`
}

func (o AppNameRule) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AppNameRule struct{}"
	}

	return strings.Join([]string{"AppNameRule", string(data)}, " ")
}
