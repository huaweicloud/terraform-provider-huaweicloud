package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type RemoveExtendCidrOption struct {

	// 功能说明：移除VPC扩展网段 取值范围：该VPC已经存在的扩展网段 约束：移除扩展网段前，请先清理该VPC下对应cidr范围内的subnet；当前只支持一个一个移除
	ExtendCidrs []string `json:"extend_cidrs"`
}

func (o RemoveExtendCidrOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemoveExtendCidrOption struct{}"
	}

	return strings.Join([]string{"RemoveExtendCidrOption", string(data)}, " ")
}
