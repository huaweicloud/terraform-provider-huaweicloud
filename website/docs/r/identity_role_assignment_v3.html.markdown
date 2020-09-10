---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role_assignment_v3"
sidebar_current: "docs-huaweicloud-resource-identity-role-assignment-v3"
description: |-
  Manages a V3 Policy assignment within HuaweiCloud IAM Service.
---

# huaweicloud\_identity\_role\_assignment_v3

Manages a V3 Role assignment within group on HuaweiCloud IAM Service.

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource. 

## Example Usage: Assign Role On Project Level

```hcl
resource "huaweicloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "huaweicloud_identity_role_v3" "role_1" {
  name = "system_all_4" #ECS admin
}

resource "huaweicloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id   = "${huaweicloud_identity_group_v3.group_1.id}"
  project_id = "${var.project_id}"
  role_id    = "${data.huaweicloud_identity_role_v3.role_1.id}"
}
```

## Example Usage: Assign Role On Domain Level

```hcl

variable "domain_id" {
  default     = "01aafcf63744d988ebef2b1e04c5c34"
  description = "this is the domain id"
}

resource "huaweicloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "huaweicloud_identity_role_v3" "role_1" {
  name = "secu_admin" #security admin
}

resource "huaweicloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = "${huaweicloud_identity_group_v3.group_1.id}"
  domain_id = "${var.domain_id}"
  role_id   = "${data.huaweicloud_identity_role_v3.role_1.id}"
}

```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `group_id` - (Required) The group to assign the role to.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.

* `role_id` - (Required) The role to assign.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `role_id` - See Argument Reference above.
