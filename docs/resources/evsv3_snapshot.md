---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_snapshot"
description: |-
  Manages an EVS snapshot (V3) resource within HuaweiCloud.
---

# huaweicloud_evsv3_snapshot

Manages an EVS snapshot (V3) resource within HuaweiCloud.

## Example Usage

```hcl
variable "volume_id" {}
variable "name" {}

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id = var.volume_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the ID of the source cloud disk for the snapshot.

* `name` - (Required, String) Specifies the name of the snapshot. Supports a maximum of `64` characters.

* `metadata` - (Optional, Map, NonUpdatable) Specifies the user-defined metadata key-value pair.

* `description` - (Optional, String) Specifies the description of the snapshot. Supports a maximum of `85` characters.

* `force` - (Optional, Bool) Specifies the flag for forcibly creating a snapshot. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the snapshot.

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

* `size` - The size of the snapshot in GiB.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

The EVS v3 snapshot can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_evsv3_snapshot.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `metadata`, `force`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the resource. Also, you can
ignore changes as below.

```hcl
resource "huaweicloud_evsv3_snapshot" "test" {
    ...

  lifecycle {
    ignore_changes = [
      metadata, force,
    ]
  }
}
```
