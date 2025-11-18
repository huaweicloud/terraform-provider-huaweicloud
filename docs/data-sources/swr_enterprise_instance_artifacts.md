---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifacts"
description: |-
  Use this data source to get the list of SWR enterprise instance artifacts.
---

# huaweicloud_swr_enterprise_instance_artifacts

Use this data source to get the list of SWR enterprise instance artifacts.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}

data "huaweicloud_swr_enterprise_instance_artifacts" "test" {
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

* `type` - (Optional, String) Specifies the artifact type.

* `tags` - (Optional, String) Specifies the artifact version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `artifacts` - Indicates the artifacts.

  The [artifacts](#artifacts_struct) structure is documented below.

* `total` - Indicates the total counts of artifact.

<a name="artifacts_struct"></a>
The `artifacts` block supports:

* `id` - Indicates the artifact ID.

* `namespace_id` - Indicates the namespace ID.

* `repository_id` - Indicates the repository ID.

* `repository_name` - Indicates the repository name.

* `size` - Indicates the artifact size, unit is byte.

* `digest` - Indicates the abstract.

* `type` - Indicates the artifact type.

* `media_type` - Indicates the media type.

* `manifest_media_type` - Indicates the manifest media type.

* `pull_time` - Indicates the last pull time.

* `push_time` - Indicates the last push time.

* `tags` - Indicates the artifact version tags.

  The [tags](#artifacts_tags_struct) structure is documented below.

<a name="artifacts_tags_struct"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `repository_id` - Indicates the repository ID.

* `artifact_id` - Indicates the artifact ID.

* `name` - Indicates the tag name.

* `push_time` - Indicates the push time of the tag.

* `pull_time` - Indicates the pull time of the tag.
