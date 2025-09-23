---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_snapshots"
description: |-
  Use this data source to get the list of EVS v3 snapshots within HuaweiCloud.
---

# huaweicloud_evsv3_snapshots

Use this data source to get the list of EVS v3 snapshots within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evsv3_snapshots" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `volume_id` - (Optional, String) Specifies the disk ID corresponding to the snapshots.

* `name` - (Optional, String) Specifies the name of the snapshots. Supports a maximum of `255` characters. This field
  will undergo a fuzzy matching query, the query result is for all snapshots whose names contain this value.

* `status` - (Optional, String) Specifies the status of the snapshots.  
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `snapshots` - The list of the snapshots.
  The [snapshots](#v3_snapshots_struct) structure is documented below.

<a name="v3_snapshots_struct"></a>
The `snapshots` block supports:

* `id` - The ID of the snapshot.

* `name` - The name of the snapshot.

* `description` - The description of the snapshot.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

* `metadata` - The user-defined metadata key-value pair.

* `volume_id` - The ID of the disk to which the snapshot belongs.

* `size` - The size of the snapshot in GiB.

* `status` - The status of the snapshot.
