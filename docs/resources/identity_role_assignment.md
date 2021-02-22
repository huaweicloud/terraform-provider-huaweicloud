---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_role\_assignment

Manages a Role assignment within group on HuaweiCloud IAM Service. This is an alternative to `huaweicloud_identity_role_assignment_v3`

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource. 

## Example Usage: Assign Role On Project Level

```hcl
resource "huaweicloud_identity_group" "group_1" {
  name = "group_1"
}

data "huaweicloud_identity_role" "role_1" {
  name = "system_all_4" #ECS admin
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  group_id   = huaweicloud_identity_group.group_1.id
  project_id = var.project_id
  role_id    = data.huaweicloud_identity_role.role_1.id
}
```

## Example Usage: Assign Role On Domain Level

```hcl

variable "domain_id" {
  default     = "01aafcf63744d988ebef2b1e04c5c34"
  description = "this is the domain id"
}

resource "huaweicloud_identity_group" "group_1" {
  name = "group_1"
}

data "huaweicloud_identity_role" "role_1" {
  name = "secu_admin" #security admin
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  group_id  = huaweicloud_identity_group.group_1.id
  domain_id = var.domain_id
  role_id   = data.huaweicloud_identity_role.role_1.id
}

```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required, String, ForceNew) The role to assign.

* `group_id` - (Required, String, ForceNew) The group to assign the role to.

* `domain_id` - (Optional, String, ForceNew; Required if `project_id` is empty) The domain to assign the role in.

* `project_id` - (Optional, String, ForceNew; Required if `domain_id` is empty) The project to assign the role in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

