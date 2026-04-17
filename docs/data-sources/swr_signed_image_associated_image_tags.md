---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_signed_image_associated_image_tags"
description: |-
  Use this data source to list the tags refer to a signed image within HuaweiCloud.
---

# huaweicloud_swr_signed_image_associated_image_tags

Use this data source to list the tags refer to a signed image within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "repository" {}
variable "sig_tag" {}

data "huaweicloud_swr_signed_image_associated_image_tags" "test" {
  namespace  = var.namespace
  repository = var.repository
  sig_tag    = var.sig_tag
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace (organization) name of the signed image belongs.

* `repository` - (Required, String) Specifies the repository name of the signed image.

* `sig_tag` - (Required, String) Specifies the tag of the signed image. Valid format is `sha256-xxx.sig`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - All tags refer to the signed image.
