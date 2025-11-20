---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_all_artifacts"
description: |-
  Use this data source to get the list of SWR enterprise artifacts within an instance.
---

# huaweicloud_swr_enterprise_instance_all_artifacts

Use this data source to get the list of SWR enterprise artifacts within an instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_instance_all_artifacts" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `artifacts` - Indicates the artifacts.

  The [artifacts](#artifacts_struct) structure is documented below.

<a name="artifacts_struct"></a>
The `artifacts` block supports:

* `id` - Indicates the artifact ID.

* `namespace_id` - Indicates the namespace ID.

* `repository_id` - Indicates the repository ID.

* `repository_name` - Indicates the repository name.

* `media_type` - Indicates the media type.

* `size` - Indicates the artifact size, unit is byte.

* `digest` - Indicates the digest.

* `type` - Indicates the artifact type.

* `manifest_media_type` - Indicates the manifest media type.

* `pull_time` - Indicates the last pull time.

* `push_time` - Indicates the last push time.

* `tags` - Indicates the artifact version tags.

  The [tags](#artifacts_tags_struct) structure is documented below.

<a name="artifacts_tags_struct"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `name` - Indicates the tag name.

* `repository_id` - Indicates the repository ID.

* `artifact_id` - Indicates the artifact ID.

* `push_time` - Indicates the push time of the tag.

* `pull_time` - Indicates the pull time of the tag.
