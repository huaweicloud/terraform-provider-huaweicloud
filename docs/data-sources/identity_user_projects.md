---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_projects"
description: |- 
  Use this data source to query the project list of a specified IAM user by administrators, or query your own project
  list within HuaweiCloud.
---

# huaweicloud_identity_user_projects

Use this data source to query the project list of a specified IAM user by administrators, or query your own project list
within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_user_projects" "test" {}
```

```hcl
data "huaweicloud_identity_user_projects" "test" {
  user_id = "The user id under your account"
}
```

## Argument Reference

* `user_id` - (Optional, String) Specifies the user id.

## Attribute Reference

* `projects` - Indicates the details of the projects.
  The [projects](#IdentityUserProjects_Projects) structure is documented below.

<a name="IdentityUserProjects_Projects"></a>
The `projects` block supports:

* `id` - Indicates the IAM project ID.

* `name` - Indicates the IAM project name.

* `enabled` - Indicates whether the IAM project is enabled.
