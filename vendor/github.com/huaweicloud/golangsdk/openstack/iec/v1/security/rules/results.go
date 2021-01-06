package rules

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

//RespSecurityGroupRule 获取安全组规则详情的结构体
type RespSecurityGroupRule struct {
	SecurityGroupRule common.RespSecurityGroupRuleEntity `json:"security_group_rule"`
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*RespSecurityGroupRule, error) {
	var entity RespSecurityGroupRule
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*RespSecurityGroupRule, error) {
	var entity RespSecurityGroupRule
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
