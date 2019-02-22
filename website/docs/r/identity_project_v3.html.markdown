---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_project_v3"
sidebar_current: "docs-huaweicloud-resource-identity-project-v3"
description: |-
  Manages a Project resource within HuaweiCloud Keystone.
---

# huaweicloud\_identity\_project_v3

Manages a Project resource within HuaweiCloud Identity And Access 
Management service.

Note: You _must_ have security admin privileges in your HuaweiCloud 
cloud to use this resource. please refer to [User Management Model](
https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html)

## Example Usage

```hcl
resource "huaweicloud_identity_project_v3" "project_1" {
  name = "eu-de_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project. it must start with 
    ID of an existing region_ and be less than or equal to 64 characters.
    Example: eu-de_project1.

* `description` - (Optional) A description of the project.

* `domain_id` - (Optional) The domain this project belongs to. Changing this
    creates a new Project.

* `parent_id` - (Optional) The parent of this project. Changing this creates
    a new Project.

* `region` - (Optional) The region in which to obtain the IAM client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new Project.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `parent_id` - See Argument Reference above.

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
