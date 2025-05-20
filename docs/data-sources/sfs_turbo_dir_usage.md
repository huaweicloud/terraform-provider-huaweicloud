---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_dir_usage"
description: |-
  Use this data source to get the usage of a directory.
---

# huaweicloud_sfs_turbo_dir_usage

Use this data source to get the usage of a directory.

-> 1. This datasource is only available for the following SFS Turbo file system types:
  **Standard**, **Standard-Enhanced**, **Performance**,**Performance-Enhanced**.
<br/>2. The obtained data may not be the latest as there is a `5` minutes delay between the frontend and background.

## Example Usage

```hcl
variable "share_id" {}
variable "path" {}

data "huaweicloud_sfs_turbo_dir_usage" "test" {
  share_id = var.share_id
  path     = var.path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS Turbo file system.

* `path` - (Required, String) Specifies the valid full path of a directory in the SFS Turbo file system.
  The `path` starts with a slash(/). e.g. **/test**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dir_usage` - The usage of the directory.

  The [dir_usage](#dir_usage_struct) structure is documented below.

<a name="dir_usage_struct"></a>
The `dir_usage` block supports:

* `used_capacity` - The used capacity of the directory, in bytes.
