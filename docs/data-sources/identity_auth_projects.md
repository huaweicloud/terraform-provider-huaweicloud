---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_auth_projects"
description: |-
  Use this data source to query the list of projects accessible by current user within HuaweiCloud.
---

# huaweicloud_identity_auth_projects

Use this data source to query the list of projects accessible by current user within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_auth_projects" "test" {}
```

## Argument Reference

This data source doesn't need argument.

## Attribute Reference

* `projects` - Indicates the details of the projects.
  The [projects](#IdentityAuthProjects_Projects) structure is documented below.

<a name="IdentityAuthProjects_Projects"></a>
The `projects` block supports:

* `id` - Indicates the IAM project ID.

* `name` - Indicates the IAM project name.

* `enabled` - Indicates whether the IAM project is enabled.
