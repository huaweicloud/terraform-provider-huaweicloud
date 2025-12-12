---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_enterprise_project_users"
description: |-
  Use this data source to get IAM users of the specified enterprise project id within HuaweiCloud.
---

# huaweicloud_identity_enterprise_project_users

Use this data source to get IAM users of the specified enterprise project id within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}

data "huaweicloud_identity_enterprise_project_users" "users" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

* `enterprise_project_id` - (Required, String) Specifies the ID of the enterprise project.

## Attribute Reference

* `users` - Indicates the user information.
  The [users](#IdentityEnterpriseProjects_Users) structure is documented below.

<a name="IdentityEnterpriseProjects_Users"></a>
The `users` block contains:

* `enabled` - Indicates whether to enable authorized users.

* `description` - Indicates the description of the IAM user.

* `domain_id` - Indicates the domain ID of the IAM user.

* `id` - Indicates the ID of the IAM user.

* `name` - Indicates the name of the IAM user.

* `policy_num` - Indicates number of user strategies

* `lastest_policy_time` - Indicates the time when the user was recently associated with the
  enterprise project strategy.

* `enterprise_projects` - Indicates the projects contains.
  The [enterprise_projects](#IdentityEnterpriseProjects_projects) structure is documented below.

<a name="IdentityEnterpriseProjects_projects"></a>
The `enterprise_projects` block contains:

* `project_id` - Indicates the ID of enterprise projects.
