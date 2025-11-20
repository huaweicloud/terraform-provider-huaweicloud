---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_repository_tags"
description: |-
  Use this data source to get the list of SWR enterprise repository tags.
---

# huaweicloud_swr_enterprise_repository_tags

Use this data source to get the list of SWR enterprise repository tags.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}

data "huaweicloud_swr_enterprise_repository_tags" "test" {
  instance_id     = var.instance_id
  namespace_name  = var.namespace_name
  repository_name = var.repository_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `repository_name` - (Required, String) Specifies the repository name.

* `is_accessory` - (Optional, String) Specifies whether to return the accessories. Value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the tags of the artifact.

  The [tags](#tags_struct) structure is documented below.

* `total` - Indicates the total counts.

<a name="tags_struct"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `name` - Indicates the tag name.

* `namespace_id` - Indicates the namespace ID.

* `repository_id` - Indicates the repository ID.

* `size` - Indicates the accessory size.

* `digest` - Indicates the digest of the accessory.

* `type` - Indicates the accessory type.

* `media_type` - Indicates the media type.

* `manifest_media_type` - Indicates the manifest media type.

* `artifact_id` - Indicates the artifact ID associated with.

* `pull_time` - Indicates the pull time of the tag.

* `push_time` - Indicates the push time of the tag.
