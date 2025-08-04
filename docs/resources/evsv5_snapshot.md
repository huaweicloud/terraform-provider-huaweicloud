---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv5_snapshot"
description: |-
  Manages an EVS snapshot (V5) resource within HuaweiCloud.
---

# huaweicloud_evsv5_snapshot

Manages an EVS snapshot (V5) resource within HuaweiCloud.

-> Before using this resource, ensure that there is no snapshot being created under the volume.

## Example Usage

```hcl
variable "volume_id" {}
variable "name" {}
variable "description" {}
variable "enterprise_project_id" {}
variable "instant_access" {}
variable "incremental" {}

resource "huaweicloud_evsv5_snapshot" "test" {
  volume_id             = var.volume_id
  name                  = var.name
  description           = var.description
  enterprise_project_id = var.enterprise_project_id
  instant_access        = var.instant_access
  incremental           = var.incremental

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the ID of the source cloud disk for the snapshot.

* `name` - (Optional, String) Specifies the name of the snapshot. Supports a maximum of `64` characters.

* `description` - (Optional, String) Specifies the description of the snapshot. Supports a maximum of `85` characters.

* `tags` - (Optional, Map) Specifies the key/value pairs to be associated with the snapshot.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID for the snapshot.

* `instant_access` - (Optional, Bool) Specifies whether to enable instant access for the snapshot. Possible values are
  **true** (enable) and **false** (disable). Default is **false**. The function not supported if the volume type is SAS
  or SATA. Only can be set **true** if it was **true** at created.

* `incremental` - (Optional, Bool, NonUpdatable) Specifies whether to create an incremental snapshot. Default is **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

* `size` - The size of the snapshot, in GiB.

* `status` - The status of the snapshot.
  The valid values are as follows:
  + **creating**: The snapshot is in the process of being created.
  + **available**: Snapshot created successfully, can be used.
  + **error**: An error occurred during the snapshot creation process.
  + **deleting**: The snapshot is in the process of being deleted.
  + **error_deleting**: An error occurred during the deletion process of the snapshot.
  + **rollbacking**: The snapshot is in the process of rolling back data.
  + **backing-up**: Through the OpenStack native API, backups can be created directly from snapshots, at this time, the
    snapshot status will change to **backing-up**. During the process of creating a backup of a disk through an API, the
    system will automatically create a snapshot, at this time, the snapshot status is **backing-up**.

* `encrypted` - Whether the snapshot is encrypted.

* `cmk_id` - The key ID of the volume to which the snapshot belongs.

* `category` - The category of snapshot.

* `availability_zone` - The AZ to which the snapshot belongs.

* `retention_at` - The duration which the snapshot is retentained.

* `instant_access_retention_at` - The retention time of the snapshot high-speed availability function. After the
  retention time expires, the snapshot high-speed availability function will be automatically disabled.

* `snapshot_type` - The snapshot creation source. Possible values as follows:
  + **auto**: The snapshot created automatically.
  + **user**: The snapshot created manually.
  + **copy**: The snapshot created by copy.

* `progress` - The snapshot creation progress.

* `encrypt_algorithm` - The algorithm of the encrypted snapshot.

* `snapshot_chains` - The snapshot chain list to which the snapshot belongs.
  The [snapshot_chains](#snapshot_chains_structure) structure is documented below.

* `snapshot_group_id` - The snapshot group ID to which the snapshot belongs.

<a name="snapshot_chains_structure"></a>
The `snapshot_chains` block supports:

* `id` - The snapshot chain ID.

* `availability_zone` - The AZ to which the snapshot chain belongs.

* `snapshot_count` - The number of snapshots on the snapshot chain.

* `capacity` - The total size of the snapshot chain.

* `volume_id` - The ID of disk to which the snapshot chain belongs.

* `category` - The category of snapshot chain.

* `created_at` - The time when the snapshot chain was created.

* `updated_at` - The time when the snapshot chain was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

The EVS v5 snapshot can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evsv5_snapshot.test <id>
```
