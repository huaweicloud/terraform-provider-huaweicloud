---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role_assignments"
description: |-
  Use this data source to get role assignments of an IAM user within HuaweiCloud.
---

# huaweicloud_identity_role_assignments

Use this data source to get role assignments of an IAM user within HuaweiCloud.

## Example Usage

### Query By DomainId

```hcl
data "huaweicloud_identity_role_assignments" "test" {}
```

### Query With Inherited

```hcl
variable "scope" {}

data "huaweicloud_identity_role_assignments" "test" {
  is_inherited = true
  scope        = var.scope
}
```

### Query By RoleId

```hcl
variable "role_id" {}
variable "scope" {}

data "huaweicloud_identity_role_assignments" "test" {
  is_inherited = true
  scope        = var.scope
  role_id      = var.role_id
}
```

## Argument Reference

* `role_id` - (Optional, String) Specifies the role ID.

* `subject` - (Optional, String) Specifies authorization entity, value range: **user**, **group**, **agency**.

* `subject_user_id` - (Optional, String) Specifies authorized IAM user ID. Conflict with `subject_agency_id` and `subject_group_id`.

* `subject_group_id` - (Optional, String) Specifies authorized user group ID. Conflict with `subject_user_id` and `subject_agency_id`.

* `subject_agency_id` - (Optional, String) Specifies authorized delegation ID. Conflict with `subject_user_id` and `subject_group_id`.

* `scope` - (Optional, String) Specifies scope of authorization, value range: **project**, **domain**, **enterprise_project**.
  Default: **domain**.

* `scope_project_id` - (Optional, String) Specifies authorized project ID. Conflict with `scope_domain_id` and
  `scope_enterprise_projects_id`.

* `scope_domain_id` - (Optional, String) Specifies pending query account ID. Conflict with `scope_enterprise_projects_id`
  and `scope_project_id`.

* `scope_enterprise_projects_id` - (Optional, String) Specifies authorized enterprise project ID. Conflict with
  `scope_domain_id` and `scope_project_id`.

* `is_inherited` - (Optional, String) Specifies takes effect when the parameter scope=domain or scope.domain_id exists.
  + **true**: query records based on authorization for all projects.
  + **false**: query records based on global service authorization.

  Default: **false**.

* `include_group` - (Optional, String) Specifies Whether to include records based on the permissions granted to the user
  groups of the IAM user.
  + **true**: Query records based on IAM user permissions and permissions granted to user groups of IAM users.
  + **false**: Only query records based on IAM user permissions.

  Default: **true**

## Attribute Reference

* `role_assignments` - Indicates authorization information.
  The [role_assignments](#IdentityRoleAssignments_Assignments) structure is documented below.

<a name="IdentityRoleAssignments_Assignments"></a>
The `role_assignments` block contains:

* `user_id` - Indicates IAM user ID.

* `role_id` - Indicates permission ID.

* `group_id` - Indicates user Group ID.

* `agency_id` - Indicates commission ID.

* `is_inherited` - Indicates  based on authorization for all projects.

* `scope` - Indicates the scope of the authorization
  The [scope](#IdentityRoleAssignments_Scope ) structure is documented below.

<a name="IdentityRoleAssignments_Scope"></a>
The `scope` block contains:

* `project_id` - Indicates information based on IAM project authorization.

* `domain_id` - Indicates information based on global service or authorization of all items.

* `enterprise_project_id` - Indicates information based on enterprise project authorization.
