---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_role_assignment"
description: |-
  Using this resource to authorize a role to a user group and specify the scope of effect within HuaweiCloud.
---

# huaweicloud_identity_group_role_assignment

Using this resource to authorize a role to a user group and specify the scope of effect within HuaweiCloud.

-> You **must** have admin privileges to use this resource.<br>
   When the resource is created, the permissions will take effect after `15` to `30` minutes.

## Example Usage

### Assign role with a specified project

```hcl
variable "group_id" {}
variable "role_name" {} # The role name can be the system name (e.g. RDS Administrator) or a custom name
variable "project_id" {}

data "huaweicloud_identity_role" "test" {
  name = var.role_name
}

resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id   = var.group_id
  role_id    = data.huaweicloud_identity_role.test.id
  project_id = var.project_id
}
```

### Assign role with all projects

```hcl
variable "group_id" {}
variable "role_name" {} # The role name can be the system name (e.g. RDS Administrator) or a custom name

data "huaweicloud_identity_role" "test" {
  name = var.role_name
}

resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id   = var.group_id
  role_id    = data.huaweicloud_identity_role.test.id
  project_id = "all"
}
```

### Assign role with a specified domain

```hcl
variable "group_id" {}
variable "role_name" {} # The role name can be the system name (e.g. RDS Administrator) or a custom name
variable "domain_id" {}

data "huaweicloud_identity_role" "test" {
  name = var.role_name
}

resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id   = var.group_id
  role_id   = data.huaweicloud_identity_role.test.id
  domain_id = var.domain_id
}
```

### Assign role with a specified enterprise project

```hcl
variable "group_id" {}
variable "role_name" {} # The role name can be the system name (e.g. RDS Administrator) or a custom name
variable "enterprise_project_id" {}

data "huaweicloud_identity_role" "test" {
  name = var.role_name
}

resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id              = var.group_id
  role_id               = data.huaweicloud_identity_role.test.id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, Nonupdatable) Specifies the ID of user group to which the role to be authorized
  belongs.

* `role_id` - (Required, String, Nonupdatable) Specifies the ID of role to be authorized.

* `domain_id` - (Optional, String, Nonupdatable) Specifies the domain to assign the role in.

* `project_id` - (Optional, String, Nonupdatable) Specifies the project to assign the role in.
  If `project_id` is set to **all**, it means that the specified user group will be able to use all projects,
  including existing and future projects.

* `enterprise_project_id` - (Optional, String, Nonupdatable) Specifies the enterprise project to assign the role in.

-> Exactly one of `domain_id`, `project_id` and `enterprise_project_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.  
  + When the role is assigned in domain, the format is `<group_id>/<role_id>/<domain_id>`
  + when the role is assigned in project, the format is `<group_id>/<role_id>/<project_id>` or
    `<group_id>/<role_id>/all`
  + when the role is assigned in enterprise project, the format is `<group_id>/<role_id>/<enterprise_project_id>`

## Import

The role assignments can be imported using the `group_id`, `role_id`, `assignment_type` and one of `domain_id`,
`project_id` and `enterprise_project_id` (assigned object ID), which depends on the value of assignment type, e.g.

```bash
$ terraform import huaweicloud_identity_group_role_assignment.test <group_id>/<role_id>/<assigned_object_id>:<assignment_type>
```

The valid values of 'assignment_type' are as follows (and the corresponding meanings of 'assigned_object_id' are also
explained):

* **domain**: the value of 'assigned_object_id' input is domain (account) ID, e.g.

```bash
$ terraform import huaweicloud_identity_group_role_assignment.test <group_id>/<role_id>/<domain_id>:domain
```

* **project**: the value of 'assigned_object_id' input is project ID or **all**, e.g.

```bash
$ terraform import huaweicloud_identity_group_role_assignment.test <group_id>/<role_id>/all:project
$ terraform import huaweicloud_identity_group_role_assignment.test <group_id>/<role_id>/<project_id>:project
```

* **enterprise_project**: the value of 'assigned_object_id' input is enterprise project ID, e.g.

```bash
$ terraform import huaweicloud_identity_group_role_assignment.test <group_id>/<role_id>/<enterprise_project_id>:enterprise_project
```
