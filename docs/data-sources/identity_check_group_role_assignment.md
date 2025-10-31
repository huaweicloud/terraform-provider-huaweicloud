---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_check_group_role_assignment"
description: |-
    Use this data source to query an IAM user group role assignment within HuaweiCloud.
---

# huaweicloud_identity_check_group_role_assignment

Use this data source to query an IAM user group role assignment within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

### Check whether the user group has global service permissions

```hcl
variable "group_id" {}
variable "role_id" {}
variable "domain_id" {}

data "huaweicloud_identity_check_group_role_assignment" "test" {
  group_id  = var.group_id
  role_id   = var.role_id
  domain_id = var.domain_id
}
```

### Check whether the user group has project service permissions

```hcl
variable "group_id" {}
variable "role_id" {}
variable "project_id" {}

data "huaweicloud_identity_check_group_role_assignment" "test" {
  group_id   = var.group_id
  role_id    = var.role_id
  project_id = var.project_id
}
```

### Check whether the user group has the specified permissions for all projects

```hcl
variable "group_id" {}
variable "role_id" {}

data "huaweicloud_identity_check_group_role_assignment" "test" {
  group_id   = var.group_id
  role_id    = var.role_id
  project_id = "all"
}
```

## Argument Reference

* `group_id` - (Required, String) Specifies the group id.

* `role_id` - (Required, String) Specifies the role id.

* `domain_id` - (Optional, String) Specifies the domain id.

* `project_id` - (Optional, String) Specifies the project id.
  If `project_id` is set to **all**, it means to check whether the user group has the specified permissions for all
  projects, including existing and future projects.

Exactly one of `domain_id` or `project_id` must be specified.

## Attribute Reference

* `result` - Indicates whether the group has the role assignment in the scope.
