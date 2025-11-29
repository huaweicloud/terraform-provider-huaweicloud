---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swrv3_image_tags"
description: |-
  Use this data source to get the list of SWR image tags.
---

# huaweicloud_swrv3_image_tags

Use this data source to get the list of SWR image tags.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swrv3_image_tags" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization.

* `repository` - (Required, String) Specifies the name of the repository.

* `tag` - (Optional, String) Specifies the source tag.

* `with_manifest` - (Optional, Bool) Specifies whether to get the manifest infos.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `schema` - Indicates the docker schema.

* `size` - Indicates the image size.

* `created` - Indicates the create time.

* `updated` - Indicates the last update time.

* `tag` - Indicates the image tag.

* `digest` - Indicates the image digest.

* `path` - Indicates the image pull path.

* `internal_path` - Indicates the image internal pull path.

* `repo_id` - Indicates the repository ID.

* `image_id` - Indicates the image ID.

* `manifest` - Indicates the image manifest.

* `is_trusted` - Indicates whether the image is trusted.

* `tag_type` - Indicates the tag type.
