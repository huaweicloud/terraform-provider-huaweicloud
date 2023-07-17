package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnterpriseProjectName 企业项目名称
type EnterpriseProjectName struct {
}

func (o EnterpriseProjectName) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnterpriseProjectName struct{}"
	}

	return strings.Join([]string{"EnterpriseProjectName", string(data)}, " ")
}
