---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_obs_target"
description: |-
  Manages an OBS target resource under the SFS Turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_obs_target

Manages an OBS target resource under the SFS Turbo within HuaweiCloud.

-> The resource supported SFS Turbo file system types are **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,
  **250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC**.

-> Due to the inherent reason of the API, updating the `attributes` may fail.

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

  lifecycle {
    ignore_changes = [
      delete_data_in_file_system, obs.0.attributes.0.file_mode, obs.0.attributes.0.dir_mode,
      obs.0.attributes.0.uid, obs.0.attributes.0.gid,
    ]
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

* `obs` - (Required, List) Specifies the detail of the OBS bucket. Changing this will create a new resource.
  The [obs](#target_obs) structure is documented below.

* `delete_data_in_file_system` - (Optional, Bool, ForceNew) Specifies whether to delete the associated directory and
  its data files in the  SFS Turbo file system when the OBS target is deleted. The default value is **false**.

<a name="target_obs"></a>
The `obs` block supports:

* `bucket` - (Required, String, ForceNew) Specifies the name of the OBS bucket.

  -> Before configuring OBS linkage, please configure bucket policies on the access control page of the OBS bucket and
    set bucket policies for sub users who need to access the OBS bucket: current bucket, all objects in the bucket,
    all operations.

* `endpoint` - (Required, String, ForceNew) Specifies the domain name of the region where the OBS bucket belongs.

* `policy` - (Optional, List) Specifies the auto synchronization policy of the storage backend.
  The [policy](#obs_policy) structure is documented below.

* `attributes` - (Optional, List) Specifies the attributes of the storage backend.
  The paramater is not supported for the file systems which are created on or before June 30, 2024 and not upgraded.
  Please submit a service ticket if you need it. [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html)
  The [attributes](#obs_attributes) structure is documented below.

<a name="obs_policy"></a>
The `policy` block supports:

* `auto_export_policy` - (Optional, List) Specifies the auto export policy of the storage backend.
  If enabled, all update made on the file system will be automatically exported to the OBS bucket.
  The [auto_export_policy](#obs_export_policy) structure is documented below.

<a name="obs_export_policy"></a>
The `attributes` block supports:

* `events` - (Optional, List) Specifies the type of the data automatically exported to the OBS bucket.
  The valid values are as follows:
  + **NEW**: Indicate add new data. Files created and then modified in the SFS Turbo interworking directory. Any data
  or metadata modifications made will be automatically synchronized to the OBS bucket.
  + **CHANGED**: Indicate modify data. Files previously imported from the OBS bucket and then modified in the SFS Turbo
  interworking directory. Any data or metadata modifications made will be automatically synchronized to the OBS bucket.
  + **DELETED**: Indicate delete data. Files deleted from the SFS Turbo interworking directory. Deletions will be
  automatically synchronized to the OBS bucket, and only such files that were previously exported to the bucket will be
  deleted.

* `prefix` - (Optional, String) Specifies the prefix to be matched in the storage backend.

* `suffix` - (Optional, String) Specifies the suffix to be matched in the storage backend.

<a name="obs_attributes"></a>
The `auto_export_policy` block supports:

* `file_mode` - (Optional, String) Specifies the permissions on the imported file.
  The valid value ranges from `0` to `777`.

* `dir_mode` - (Optional, String) Specifies the permissions on the imported directory.
  The valid value ranges from `0` to `777`.

-> For more details about the fields, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-sfsturbo/CreateBackendTarget.html).

* `uid` - (Optional, Int) Specifies the ID of the user who owns the imported object. Default value is `0`.
  The valid value ranges from `0` to `4,294,967,294`.

* `gid` - (Optional, Int) Specifies the ID of the user group to which the imported object belongs.
  Default value is `0`. The valid value ranges from `0` to `4,294,967,294`.

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
The missing attributes include: `delete_data_in_file_system`, `obs.0.attributes.0.file_mode`,
`obs.0.attributes.0.dir_mode`, `obs.0.attributes.0.uid`, `obs.0.attributes.0.gid`.
It is generally recommended running `terraform plan` after importing an resource.
You can ignore changes as below.

```hcl
resource "huaweicloud_sfs_turbo_obs_target" "test" {
  ...

  lifecycle {
    ignore_changes = [
      delete_data_in_file_system, obs.0.attributes.0.file_mode, obs.0.attributes.0.dir_mode,
      obs.0.attributes.0.uid, obs.0.attributes.0.gid,
    ]
  }
}
```
