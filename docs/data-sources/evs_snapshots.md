---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshots"
description: ""
---

# huaweicloud_evs_snapshots

Use this data source to get a list of EVS cloudvolume snapshots.

## Example Usage

```hcl
variable "volume_id" {}

data "huaweicloud_evs_snapshots" "snapshots" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the WAF dedicated domains.
  If omitted, the provider-level region will be used.
  
* `enterprise_project_id` - (Optional, String) Specify the enterprise project ID to filter.
  If this field value is not specified, cloudvolume snapshots of all enterprise projects within
  authority scope will be queried.

* `name` - (Optional, String)  The name of cloudvolume snapshot. Maximum supported is 64 characters.

* `status` - (Optional, String) The status of cloudvolume snapshot. Valid values are as follows:
  + **creating**: The cloudvolume snapshot is in the process of being created.
  + **available**: The cloudvolume snapshot is created successfully and can be used.
  + **error**: An error occurred during the creation process of the cloudvolume snapshot.
  + **deleting**: The cloudvolume snapshot is in the process of being deleted.
  + **error_deleting**: An error occurred during the deletion process of the cloudvolume snapshot.
  + **rollbacking**: The cloudvolume snapshot is in the process of rolling back data.
  + **backing-up**: The status of the following two snapshots is **backing_up**:
  The snapshots that can create backup directly through the OpenStack native API.
  The snapshots created automatically during the process of creating a backup.

* `volume_id` - (Optional, String) The ID of the cloudvolume to which the snapshot belongs.

* `availability_zone` - (Optional, String) The availability zone of the cloudvolume to which the snapshot belongs.

* `snapshot_id` - (Optional, String) Specify the snapshot ID to filter.
  You can pass in multiple ids to filter the query, for example: id=id1&id=id2&id=id3.

* `dedicated_storage_name` - (Optional, String) The name of the dedicated storage.

* `dedicated_storage_id` - (Optional, String) The ID of the dedicated storage.

* `service_type` - (Optional, String) Service type. Only **EVS**, **DSS**, and **DESS** are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `snapshots` - A list of EVS cloudvolume snapshots.

The `snapshots` block supports:

* `id` - The data source ID of EVS cloudvolume snapshot.

* `name` - The name of cloudvolume snapshot. Maximum supported is 64 characters.

* `size` - The cloudvolume snapshot size. Unit is GiB.

* `status` - The status of cloudvolume snapshot. Valid values are as follows:
  + **creating**: The cloudvolume snapshot is in the process of being created.
  + **available**: The cloudvolume snapshot is created successfully and can be used.
  + **error**: An error occurred during the creation process of the cloudvolume snapshot.
  + **deleting**: The cloudvolume snapshot is in the process of being deleted.
  + **error_deleting**: An error occurred during the deletion process of the cloudvolume snapshot.
  + **rollbacking**: The cloudvolume snapshot is in the process of rolling back data.
  + **backing-up**: The status of the following two snapshots is **backing_up**:
  The snapshots that can create backup directly through the OpenStack native API.
  The snapshots created automatically during the process of creating a backup.

* `description` - The cloudvolume snapshot description information.

* `created_at` - The cloudvolume snapshot creation time.

* `updated_at` - The cloudvolume snapshot update time.

* `volume_id` - The ID of the cloudvolume to which the snapshot belongs.

* `service_type` - Service type. Valid values can be **EVS**, **DSS**, and **DESS**.

* `metadata` - The user-defined metadata key-value pair.

* `dedicated_storage_id` - The dedicated storage pool ID.

* `dedicated_storage_name` - The dedicated storage pool name.

* `progress` - The snapshot creation progress.
