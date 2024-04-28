---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role_assignment"
description: ""
---

# huaweicloud_identity_role_assignment

Manages a Role assignment within group on HuaweiCloud IAM Service.

-> **NOTE:** 1. You *must* have admin privileges to use this resource.
  <br/>2. When the resource is created, the permissions will take effect after 15 to 30 minutes.

## Example Usage: Assign Role On Project Level

```hcl
data "huaweicloud_identity_role" "role_1" {
  # RDS Administrator
  name = "rds_adm"
}

resource "huaweicloud_identity_group" "group_1" {
  name = "group_1"
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  role_id    = data.huaweicloud_identity_role.role_1.id
  group_id   = huaweicloud_identity_group.group_1.id
  project_id = var.project_id
}
```

## Example Usage: Assign Role On Domain Level

```hcl
data "huaweicloud_identity_role" "role_1" {
  # Security Administrator
  name = "secu_admin"
}

resource "huaweicloud_identity_group" "group_1" {
  name = "group_1"
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  role_id   = data.huaweicloud_identity_role.role_1.id
  group_id  = huaweicloud_identity_group.group_1.id
  domain_id = var.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required, String, ForceNew) Specifies the role to assign.
  Changing this parameter will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the group to assign the role to.
  Changing this parameter will create a new resource.

* `domain_id` - (Optional, String, ForceNew) Specifies the domain to assign the role in.
  Required if `project_id` is empty. Changing this parameter will create a new resource.

* `project_id` - (Optional, String, ForceNew) Specifies the project to assign the role in.
  Required if `domain_id` is empty. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
