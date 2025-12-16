---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_repository_tag"
description: |-
  Manages a SWR repository tag resource within HuaweiCloud.
---

# huaweicloud_swr_repository_tag

Manages a SWR repository tag resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}
variable "source_tag" {}
variable "destination_tag" {}

resource "huaweicloud_swr_repository_tag" "test" {
  organization    = var.organization
  repository      = var.repository
  source_tag      = var.source_tag
  destination_tag = var.destination_tag
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `organization` - (Required, String, NonUpdatable) Specifies the name of the organization.

* `repository` - (Required, String, NonUpdatable) Specifies the name of the repository.

* `source_tag` - (Required, String, NonUpdatable) Specifies the source tag.

* `destination_tag` - (Required, String, NonUpdatable) Specifies the destination tag.

* `override` - (Optional, Bool, NonUpdatable) Specifies whether to override. Default to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created` - Indicates the creation time.

* `digest` - Indicates the image digest.

* `image_id` - Indicates the image ID.

* `internal_path` - Indicates the image internal pull path.

* `is_trusted` - Indicates whether the image is trusted.

* `manifest` - Indicates the image manifest.

* `path` - Indicates the image pull path.

* `repository_id` - Indicates the repository ID.

* `schema` - Indicates the docker schema.

* `size` - Indicates the image size.

* `tag_id` - Indicates the tag ID.

* `tag` - Indicates the tag.

* `tag_type` - Indicates the tag type.

* `updated` - Indicates the last update time.

## Import

The repository tag can be imported using `organization`, `repository` and `tag`, e.g.

Only when repository name is with no slashes, can use slashes to separate.

```bash
$ terraform import huaweicloud_swr_repository_tag.test <organization>/<repository>/<tag>
```

Using comma to separate is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_repository_tag.test <organization>,<repository>,<tag>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `source_tag`, `destination_tag` and `override`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_swr_repository_tag" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_tag, destination_tag, override
    ]
  }
}
```
