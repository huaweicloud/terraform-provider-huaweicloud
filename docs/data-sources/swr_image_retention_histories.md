---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_retention_histories"
description: |-
  Use this data source to get the list of SWR image retention histories.
---

# huaweicloud_swr_image_retention_histories

Use this data source to get the list of SWR image retention histories.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_image_retention_histories" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization.

* `repository` - (Required, String) Specifies the image repository name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The image retention histories.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The ID of the image retention history record.

* `retention_id` - The image retention policy ID.

* `organization` - The organization name.

* `repository` - The image repository name.

* `rule_type` - The image retention rule type.

* `tag` - The image tag.

* `created_at` - The creation time.
