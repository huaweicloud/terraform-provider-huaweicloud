package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListPermissionsRequest struct {

	// 系统内部呈现的权限名称。如云目录服务CCS普通用户权限CCS User的name为ccs_user。 建议您传参display_name，不传name参数。
	Name *string `json:"name,omitempty"`

	// 账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。 > - 如果填写此参数，则返回账号下所有自定义策略。 > - 如果不填写此参数，则返回所有系统权限（包含系统策略和系统角色）。
	DomainId *string `json:"domain_id,omitempty"`

	// 分页查询时数据的页数，查询值最小为1。需要与per_page同时存在。传入domain_id参数查询自定义策略时，可配套使用。
	Page *int32 `json:"page,omitempty"`

	// 分页查询时每页的数据个数，取值范围为[1,300]，默认值为300。需要与page同时存在。不传page和per_page参数时，每页最多返回300个权限。
	PerPage *int32 `json:"per_page,omitempty"`

	// 区分系统权限类型的参数。当domain_id参数为空时生效。 > - policy：返回系统策略。 > - role：返回系统角色。
	PermissionType *string `json:"permission_type,omitempty"`

	// 过滤权限名称。如传参为Administrator，则返回满足条件的所有管理员权限。
	DisplayName *string `json:"display_name,omitempty"`

	// 过滤权限的显示模式。取值范围：domain,project,all。type为domain时，返回type=AA或AX的权限；type为project时，返回type=AA或XA的权限；type为all时返回type为AA、AX、XA的权限。 > - AX表示在domain层显示。 > - XA表示在project层显示。 > - AA表示在domain和project层均显示。 > - XX表示在domain和project层均不显示。
	Type *string `json:"type,omitempty"`

	// 权限所在目录。catalog值精确匹配策略的catalog字段(可以过滤服务的策略、或者自定义策略)。
	Catalog *string `json:"catalog,omitempty"`
}

func (o KeystoneListPermissionsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListPermissionsRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListPermissionsRequest", string(data)}, " ")
}
