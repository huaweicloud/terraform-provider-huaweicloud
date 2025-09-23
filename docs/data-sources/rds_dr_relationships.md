---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_dr_relationships"
description: |-
  Use this data source to get the list of DR relationships.
---

# huaweicloud_rds_dr_relationships

Use this data source to get the list of DR relationships.

## Example Usage

```hcl
data "huaweicloud_rds_dr_relationships" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `relationship_id` - (Optional, String) Specifies the DR relationship ID.

* `status` - (Optional, String) Specifies the DR configuration status.

* `master_instance_id` - (Optional, String) Specifies the primary instance ID.

* `master_region` - (Optional, String) Specifies the region where the primary instance is located.

* `slave_instance_id` - (Optional, String) Specifies the DR instance ID.

* `slave_region` - (Optional, String) Specifies the region where the DR instance is located.

* `create_at_start` - (Optional, Int) Specifies the creation start time.

* `create_at_end` - (Optional, Int) Specifies the creation end time.

* `order` - (Optional, String) Specifies the sorting order. Value options:
  + **DESC**: descending order.
  + **ASC**: ascending order.

  Defaults to **DESC**.

* `sort_field` - (Optional, String) Specifies the sorting field.
  + **status**: Data is sorted by DR configuration status.
  + **time**: Data is sorted by DR configuration time.
  + **master_region**: Data is sorted by region where the primary instance is located.
  + **slave_region**: Data is sorted by region where the DR instance is located.

  Defaults to **time**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_dr_infos` - Indicates the list of DR relationships.
  The [instance_dr_infos](#instance_dr_infos_struct) structure is documented below.

<a name="instance_dr_infos_struct"></a>
The `instance_dr_infos` block supports:

* `id` - Indicates the DR relationship ID.

* `status` - Indicates the DR configuration status.

* `failed_message` - Indicates the failure message.

* `master_instance_id` - Indicates the primary instance ID.

* `master_region` - Indicates the region where the primary instance is located.

* `slave_instance_id` - Indicates the DR instance ID.

* `slave_region` - Indicates the region where the standby instance is located.

* `build_process` - Indicates the process for configuring disaster recovery (DR). The value can be:
  + **master**: process of configuring DR capability for the primary instance.
  + **slave**: process of configuring DR for the DR instance.

* `time` - Indicates the DR configuration time.

* `replica_state` - Indicates the synchronization status. The value can be:
  + **0**: indicates that the synchronization is normal.
  + **-1** indicates that the synchronization is abnormal.

* `wal_write_receive_delay_in_mb` - Indicates the WAL send lag volume, in MB. It means the difference between the WAL Log
  Sequence Number (LSN) written by the primary instance and the WAL LSN received by the DR instance.

* `wal_write_replay_delay_in_mb` - Indicates the end-to-end delayed WAL size, in MB. It refers to the difference between
  the WAL LSN written by the primary instance and the WAL LSN replayed by the DR instance.

* `wal_receive_replay_delay_in_ms` - Indicates the replay delay, in milliseconds, on the DR instance.
