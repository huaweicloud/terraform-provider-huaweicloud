---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_recycling_instances"
description: |-
  Use this data source to get the list of RDS recycling instances.
---

# huaweicloud_rds_recycling_instances

Use this data source to get the list of RDS recycling instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_recycling_instances" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the instance name.

* `ha_mode` - (Optional, String) Specifies the instance type.
  Value options: **Ha**, **Single**.

* `engine_name` - (Optional, String) Specifies the DB engine name.

* `engine_version` - (Optional, String) Specifies the DB engine version.

* `pay_model` - (Optional, String) Specifies the billing mode.
  Value options: **0** (pay-per-use), **1** (yearly/monthly).

* `volume_type` - (Optional, String) Specifies the storage type.
  Value options:
  + **ULTRAHIGH**: ultra-high I/O storage.
  + **ULTRAHIGHPRO**: ultra-high I/O (advanced) storage.
  + **CLOUDSSD**: cloud SSD storage.
  + **LOCALSSD**: local SSD storage.

* `volume_size` - (Optional, String) Specifies the storage space in **GB**.
  The value must be a multiple of **10** and the value range is from **40 GB** to **4,000 GB**.

* `data_vip` - (Optional, String) Specifies the floating IP address.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `recycle_backup_id` - (Optional, String) Specifies the backup ID.

* `recycle_status` - (Optional, String) Specifies the backup status.
  Value options:
  + **BUILDING**: The instance is being backed up and cannot be rebuilt.
  + **COMPLETED**: The backup is complete and the instance can be rebuilt.

* `is_serverless` - (Optional, String) Specifies whether the instance is a serverless instance.
  Value options: **true**, **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of recycling RDS instances.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `ha_mode` - Indicates the instance type.

* `engine_name` - Indicates the instance DB engine name.

* `engine_version` - Indicates the instance DB engine version.

* `pay_model` - Indicates the instance billing mode.

* `volume_type` - Indicates the storage type.

* `volume_size` - Indicates the storage space in **GB**.

* `data_vip` - Indicates the floating IP address.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `retained_until` - Indicates the retention time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `recycle_backup_id` - Indicates the backup ID.

* `recycle_status` - Indicates the backup status.

* `is_serverless` - Indicates whether the instance is a serverless instance.

* `created_at` - Indicates the creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `deleted_at` - Indicates the deletion time in the **yyyy-mm-ddThh:mm:ssZ** format.
