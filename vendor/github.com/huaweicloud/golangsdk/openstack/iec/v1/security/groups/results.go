package groups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

//RespSecurityGroupEntity 获取安全组信息时候，返回结构体
type RespSecurityGroupEntity struct {
	ID                   string                               `json:"id"`
	Name                 string                               `json:"name"`
	VpcID                string                               `json:"vpc_id,omitempty"`
	Description          string                               `json:"description"`
	SecurityGroupRules   []common.RespSecurityGroupRuleEntity `json:"security_group_rules"`
	RegionSecurityGroups []common.RegionSecurityGroupItem     `json:"region_security_groups,omitempty"`
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*RespSecurityGroupEntity, error) {
	var entity RespSecurityGroupEntity
	err := r.ExtractIntoStructPtr(&entity, "security_group")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*RespSecurityGroupEntity, error) {
	var entity RespSecurityGroupEntity
	err := r.ExtractIntoStructPtr(&entity, "security_group")
	return &entity, err
}
