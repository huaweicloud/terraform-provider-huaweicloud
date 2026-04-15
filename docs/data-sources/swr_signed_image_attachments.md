---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_signed_image_attachments"
description: |-
  Use this data source to list the attachments of a signed image within HuaweiCloud.
---

# huaweicloud_swr_signed_image_attachments

Use this data source to list the attachments of a signed image within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "repository" {}
variable "tag" {}

data "huaweicloud_swr_signed_image_attachments" "test" {
  namespace  = var.namespace
  repository = var.repository
  tag        = var.tag
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace (organization) name of the signed image belongs.

* `repository` - (Required, String) Specifies the repository name of the signed image.

* `tag` - (Required, String) Specifies the tag of the signed image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accessories` - All accessories that match the filter parameters.
  The [accessories](#accessories) structure is documented below.

<a name="accessories"></a>
The `accessories` block supports:

* `id` - The ID of the attachment.

* `domain_id` - The ID of the tenant the attachment belongs to.

* `namespace_name` - The name of the organization (namespace) the attachment belongs to.

* `repo_name` - The name of the repository the attachment belongs to.

* `sig_tag` - The signature tag of the attachment.

* `sig_digest` - The hash value of the attachment.

* `target_digest` - The hash value of the signed image associated with the attachment.

* `size` - The size of the attachment.

* `type` - The type of the attachment.

* `created_at` - The creation time of the attachment.

* `updated_at` - The update time of the attachment.
