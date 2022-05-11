package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AgencyTokenUserDomain struct {

	// 委托方A的账号ID。
	Id string `json:"id"`

	// 委托方A的账号名称。
	Name string `json:"name"`
}

func (o AgencyTokenUserDomain) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyTokenUserDomain struct{}"
	}

	return strings.Join([]string{"AgencyTokenUserDomain", string(data)}, " ")
}
