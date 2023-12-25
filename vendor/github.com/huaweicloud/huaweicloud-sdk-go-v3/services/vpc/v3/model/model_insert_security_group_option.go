package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InsertSecurityGroupOption
type InsertSecurityGroupOption struct {

	// 功能说明：安全组的ID列表；例如：\"security_groups\": [\"a0608cbf-d047-4f54-8b28-cd7b59853fff\"]
	SecurityGroups []string `json:"security_groups"`

	// 安全组插入的位置，从0开始计数。 举例： 1. 要插入到已关联安全组列表的首位，index=0； 2. 要插入到已关联安全组列表的第n个安全组后面，index=n。 默认插入到端口已关联的安全组列表末尾。
	Index *int32 `json:"index,omitempty"`
}

func (o InsertSecurityGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InsertSecurityGroupOption struct{}"
	}

	return strings.Join([]string{"InsertSecurityGroupOption", string(data)}, " ")
}
