package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainRoleAssignmentsRequest Request Object
type ShowDomainRoleAssignmentsRequest struct {

	// 待查询账号ID。
	DomainId string `json:"domain_id"`

	// 策略ID。
	RoleId *string `json:"role_id,omitempty"`

	// 授权主体,取值范围：user、group、agency。该参数与subject.user_id、subject.group_id、subject.agency_id只能选择一个。
	Subject *string `json:"subject,omitempty"`

	// 授权的IAM用户ID。
	SubjectUserId *string `json:"subject.user_id,omitempty"`

	// 授权的用户组ID。
	SubjectGroupId *string `json:"subject.group_id,omitempty"`

	// 授权的委托ID。
	SubjectAgencyId *string `json:"subject.agency_id,omitempty"`

	// 授权范围，取值范围：project、domain、enterprise_project。该参数与scope.project_id、scope.domain_id、scope.enterprise_projects_id只能选择一个。 > - 如需查看全局服务授权记录，scope取值domain或填写scope.domain_id。 > - 如需查看基于所有资源的授权记录，scope取值为domain，且is_inherited取值为true > - 如需查看基于项目的授权记录，scope取值为project或填写scope.project_id。 > - 如需查看基于企业项目的授权记录，scope取值为enterprise_project或填写scope.enterprise_project_id。
	Scope *string `json:"scope,omitempty"`

	// 授权的项目ID。
	ScopeProjectId *string `json:"scope.project_id,omitempty"`

	// 待查询账号ID。
	ScopeDomainId *string `json:"scope.domain_id,omitempty"`

	// 授权的企业项目ID。
	ScopeEnterpriseProjectsId *string `json:"scope.enterprise_projects_id,omitempty"`

	// 是否包含基于所有项目授权的记录，默认为false。当参数scope=domain或者scope.domain_id存在时生效。true：查询基于所有项目授权的记录。 false：查询基于全局服务授权的记录。
	IsInherited *bool `json:"is_inherited,omitempty"`

	// 是否包含基于IAM用户所属用户组授权的记录，默认为true。当参数subject=user或者subject.user_id存在时生效。true：查询基于IAM用户授权、IAM用户所属用户组授权的记录。 false：仅查询基于IAM用户授权的记录。
	IncludeGroup *bool `json:"include_group,omitempty"`

	// 分页查询时数据的页数，查询值最小为1。需要与per_page同时存在。
	Page *int32 `json:"page,omitempty"`

	// 分页查询时每页的数据个数，取值范围为[1,50]。需要与page同时存在。
	PerPage *int32 `json:"per_page,omitempty"`
}

func (o ShowDomainRoleAssignmentsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainRoleAssignmentsRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainRoleAssignmentsRequest", string(data)}, " ")
}
