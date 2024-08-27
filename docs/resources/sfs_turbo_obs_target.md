---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_obs_target"
description: ""
---

# huaweicloud_sfs_turbo_obs_target

Manages an OBS target resource under the SFS Turbo within HuaweiCloud.

-> The resource supported SFS Turbo file system types are **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,
  **250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC**.

## Example Usage

```hcl
variable "share_id" {}
variable "file_path" {}
variable "bucket_name" {}
variable "endpoint" {}

resource "huaweicloud_sfs_turbo_obs_target" "test" {
  share_id         = var.share_id
  file_system_path = var.file_path

  obs {
    bucket   = var.bucket_name
    endpoint = var.endpoint
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, ForceNew) Specifies the ID of the SFS Turbo file system to which the OBS target
  belongs. Changing this creates a new resource.

* `file_system_path` - (Required, String, ForceNew) Specifies the linkage directory name of the OBS target.
  Changing this creates a new resource.

  -> The directory name must be unique and it can not be `.` or `..` character. The directory name can not contain
    slashes (/) and multi level directory is not supported.

* `obs` - (Required, List, ForceNew) Specifies the detail of the OBS bucket. Changing this will create a new resource.
  The [obs](#target_obs) structure is documented below.

* `delete_data_in_file_system` - (Optional, Bool, ForceNew) Specifies whether to delete the associated directory and
  its data files in the  SFS Turbo file system when the OBS target is deleted. The default value is **false**.

<a name="target_obs"></a>
The `obs` block supports:

* `bucket` - (Required, String) Specifies the name of the OBS bucket.

  -> Before configuring OBS linkage, please configure bucket policies on the access control page of the OBS bucket and
    set bucket policies for sub users who need to access the OBS bucket: current bucket, all objects in the bucket,
    all operations.

* `endpoint` - (Required, String) Specifies the domain name of the region where the OBS bucket belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the OBS target. The valid values are as following:
  + **AVAILABLE**: The resource is available.
  + **MISCONFIGURED**: The resource creation failed.
  + **FAILED**: The resource deletion failed.

* `created_at` - The creation time of the OBS target.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The SFS Turbo OBS target can be imported using the related `share_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_sfs_turbo_obs_target.test <share_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to payment attributes missing from
the API response.
The missing attributes include: `delete_data_in_file_system`.
It is generally recommended running `terraform plan` after importing an resource.
You can ignore changes as below.

```hcl
resource "huaweicloud_sfs_turbo_obs_target" "test" {
  ...

  lifecycle {
    ignore_changes = [
      delete_data_in_file_system,
    ]
  }
}
```
