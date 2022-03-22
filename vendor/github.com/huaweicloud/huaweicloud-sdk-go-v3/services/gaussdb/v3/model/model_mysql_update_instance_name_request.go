package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlUpdateInstanceNameRequest struct {
	// 实例名称。 用于表示实例的名称，同一租户下，同类型的实例名可重名。取值范围：4~64个字符之间，必须以字母开头，区分大小写，可以包含字母、数字、中划线或者下划线，不能包含其他的特殊字符。

	Name string `json:"name"`
}

func (o MysqlUpdateInstanceNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlUpdateInstanceNameRequest struct{}"
	}

	return strings.Join([]string{"MysqlUpdateInstanceNameRequest", string(data)}, " ")
}
