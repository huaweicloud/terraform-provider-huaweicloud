---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_tags"
description: ""
---

# huaweicloud_swr_image_tags

Use this data source to get the list of SWR image tags.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_image_tags" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization.

* `repository` - (Required, String) Specifies the name of the repository.

* `name` - (Optional, String) Specifies the name of the image tag.

* `digest` - (Optional, String) Specify the hash value of the image tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `image_tags` - All image tags that match the filter parameters.
  The [image_tags](#attrblock--image_tags) structure is documented below.

<a name="attrblock--image_tags"></a>
The `image_tags` block supports:

* `name` - The name of the image tag.

* `size` - The size of the image tag in byte.

* `path` - The image address for docker pull.

* `internal_path` - The intra-cluster image address for docker pull.

* `digest` - The hash value of the image tag.

* `image_id` - The ID of the image.

* `is_trusted` - Whether the image version is trusted.

* `manifest` - The manifest of the image tag.

* `scanned` - Whether the image version is scanned.

* `docker_schema` - The docker protocol used by the image tag.

* `type` - The type of the image tag.

* `created_at` - The creation time of the image tag.

* `updated_at` - The update time of the image tag.

* `deleted_at` - The delete time of the image tag.
