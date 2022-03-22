package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlCreateBackupRequest struct {
	// 实例ID，严格匹配UUID规则。

	InstanceId string `json:"instance_id"`
	// 备份名称。 取值范围：4~64个字符之间，必须以字母开头，区分大小写，可以包含字母、数字、中划线或者下划线，不能包含其他的特殊字符。

	Name string `json:"name"`
	// 备份描述，不能包含>!<\"&'=特殊字符，不大于256个字符。

	Description *string `json:"description,omitempty"`
}

func (o MysqlCreateBackupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlCreateBackupRequest struct{}"
	}

	return strings.Join([]string{"MysqlCreateBackupRequest", string(data)}, " ")
}
