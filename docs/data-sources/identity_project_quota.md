---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_project_quota"
description: |-
  Use this data source to get project quota.
---

# huaweicloud_identity_project_quota

Use this data source to get project quota.

## Example Usage

```hcl
variable "project_id" {}

data "huaweicloud_identity_project_quota" "test" {
  project_id = var.project_id
}
```

## Argument Reference

* `project_id` - (Required, String) Specifies the project id.

## Attribute Reference

* `resources` - Indicates the resource info list.
  The [resources](#IdentityProjectQuota_Resources) structure is documented below.

<a name="IdentityProjectQuota_Resources"></a>
The `resources` block contains:

* `type` - Indicates the type of resource, must be **project**.

* `max` - Indicates the maximum quota.

* `min` - Indicates the minimum quota.

* `quota` - Indicates the current quota.

* `used` - Indicates the used quota.
