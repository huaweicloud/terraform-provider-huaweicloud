---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_auto_sync"
description: ""
---

# huaweicloud_swr_image_auto_sync

Manages a SWR image auto sync within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "repository_name" {}

resource "huaweicloud_swr_image_auto_sync" "test"{
  organization        = var.organization_name
  repository          = var.repository_name
  target_region       = "cn-north-4"
  target_organization = "target_org_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization.

  Changing this parameter will create a new resource.

* `repository` - (Required, String, ForceNew) Specifies the name of the repository.

  Changing this parameter will create a new resource.

* `target_region` - (Required, String, ForceNew) Specifies the target region name.

  Changing this parameter will create a new resource.

* `target_organization` - (Required, String, ForceNew) Specifies the target organization name.

  Changing this parameter will create a new resource.

* `override` - (Optional, Bool, ForceNew) Specifies whether to overwrite.
  Default to **false**, which indicates not to overwrite
  any nonidentical image that has the same name in the target organization.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

## Import

The SWR image auto sync can be imported using the organization name, repository name,
target region and target organization separated by slashes or commas, e.g.:

Only when repository name is with no slashes, can use slashes to separate.

```bash
$ terraform import huaweicloud_swr_image_auto_sync.test <organization_name>/<repository_name>/<target_region>/<target_organization>
```

Using comma to separate is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_image_auto_sync.test <organization_name>,<repository_name>,<target_region>,<target_organization>
```
