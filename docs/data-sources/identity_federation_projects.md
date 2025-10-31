---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_federation_projects"
description: |-
    Use this data source to query the list of projects accessible by federated users within HuaweiCloud.
---

# huaweicloud_identity_federation_projects

Use this data source to query the list of projects accessible by federated users within HuaweiCloud.

## Example Usage

```hcl
variable "federation_token" {}

data "huaweicloud_identity_federation_projects" "test" {
  federation_token = var.federation_token
}
```

## Argument Reference

* `federation_token` - (Required, String) Specifies federated authentication unscoped token.

## Attribute Reference

* `projects` - Indicates the details of the projects.
  The [projects](#IdentityFederationProjects_Projects) structure is documented below.

<a name="IdentityFederationProjects_Projects"></a>
The `projects` block supports:

* `id` - Indicates the IAM project ID.

* `name` - Indicates the IAM project name.

* `enabled` - Indicates whether the IAM project is enabled.
