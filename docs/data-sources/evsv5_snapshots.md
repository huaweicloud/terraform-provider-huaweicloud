---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv5_snapshots"
description: |-
  Use this data source to query the list of EVS v5 snapshots within HuaweiCloud.
---

# huaweicloud_evsv5_snapshots

Use this data source to query the list of EVS v5 snapshots within HuaweiCloud.

## Example Usage

```hcl
variable "volume_id" {}

data "huaweicloud_evsv5_snapshots" "test" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `volume_id` - (Optional, String) Specifies the ID of disk to which the snapshot belongs.

* `availability_zone` - (Optional, String) Specifies the AZ to which the snapshot chain belongs.

* `name` - (Optional, String) Specifies the snapshot name.

* `status` - (Optional, String) Specifies the snapshot status.
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

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The value can be **id**, **status**, or **created_at**, and the default value is **created_at**.

* `sort_dir` - (Optional, String) Specifies the result sorting order. The default value is **desc**.
  + **desc**: The descending order.
  + **asc**: The ascending order.

* `id` - (Optional, String) Specifies the snapshot ID.

* `ids` - (Optional, String) Specifies the snapshot IDs. The value is in the following
  format: **ids=id1,id2,...,idx**. Returns snapshot information corresponding to a valid id. Invalid ids will be ignored.

* `snapshot_type` - (Optional, String) Specifies the snapshot type.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID for filtering.

* `tag_key` - (Optional, String) Specifies the tag name used to filter results.

* `tags` - (Optional, String) Specifies the key/value pairs used to filter results. The value is in the following
  format: **{"key1":"value1"}**

* `snapshot_chain_id` - (Optional, String) Specifies the snapshot chain ID.

* `snapshot_group_id` - (Optional, String) Specifies the snapshot group ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshots` - The snapshot list.
  The [snapshots](#snapshots_structure) structure is documented below.

<a name="snapshots_structure"></a>
The `snapshots` block supports:

* `id` - The snapshot ID.

* `name` - The snapshot name. Snapshots with the prefix autobk_snapshot_ are automatically created by the system when
  you create a cloud disk backup. Do not perform the "Delete Cloud Disk Snapshot" or "Roll Back Snapshot to Cloud
  Disk" operations.

* `description` - The snapshot description.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

* `volume_id` - The ID of disk to which the snapshot belongs.

* `size` - The total size of the snapshot, in GiB.

* `status` - The snapshot status.

* `enterprise_project_id` - The ID of the enterprise project associated with the snapshot.

* `encrypted` - Whether the snapshot is encrypted.

* `cmk_id` - The custom key for the disk id to which the encrypted snapshot belongs.

* `category` - The category of snapshot.

* `availability_zone` - The AZ to which the snapshot belongs.

* `tags` - The tags list.
  The [tags](#tags_structure) structure is documented below.

* `instant_access` - Whether the snapshot high-speed availability function is enabled. Possible values as follows:
  + **true**: Enable. Only SSD series cloud disks support this function.
  + **false**: Disable. The snapshot is an existing snapshot without the snapshot high-speed availability function enabled.

* `retention_at` - The time which the snapshot ID retention at.

* `instant_access_retention_at` - The retention time of the snapshot high-speed availability function. After the
  retention time expires, the snapshot high-speed availability function will be automatically disabled.

* `incremental` - Whether the snapshot is incremental snapshot.

* `snapshot_type` - The snapshot created type. Possible values as follows:
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

<a name="tags_structure"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
