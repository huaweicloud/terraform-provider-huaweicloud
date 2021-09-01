package groups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
)

type commonResult struct {
	golangsdk.Result
}

//RespSecurityGroupEntity 获取安全组信息时候，返回结构体
type RespSecurityGroupEntity struct {
	ID                 string                               `json:"id"`
	Name               string                               `json:"name"`
	Description        string                               `json:"description"`
	SecurityGroupRules []common.RespSecurityGroupRuleEntity `json:"security_group_rules"`
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

// SecurityGroups 安全组列表对象
type SecurityGroups struct {
	SecurityGroups []RespSecurityGroupEntity `json:"security_groups"`
	Count          int                       `json:"count"`
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*SecurityGroups, error) {
	var entity SecurityGroups
	err := r.ExtractIntoStructPtr(&entity, "")
	return &entity, err
}
