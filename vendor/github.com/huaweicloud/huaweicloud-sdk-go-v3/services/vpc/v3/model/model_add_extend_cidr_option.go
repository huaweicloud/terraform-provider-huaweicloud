package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AddExtendCidrOption struct {
	// 功能说明：扩展cidr列表 取值范围：不能包含以下网段， 10.0.0.0/8     172.16.0.0/12  192.168.0.0/16 100.64.0.0/10 214.0.0.0/7 198.18.0.0/15  169.254.0.0/16 0.0.0.0/8 127.0.0.0/8  240.0.0.0/4  255.255.255.255/32  约束：当前只支持添加一个

	ExtendCidrs []string `json:"extend_cidrs"`
}

func (o AddExtendCidrOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddExtendCidrOption struct{}"
	}

	return strings.Join([]string{"AddExtendCidrOption", string(data)}, " ")
}
