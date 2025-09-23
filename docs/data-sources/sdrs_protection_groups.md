---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protection_groups"
description: |-
  Use this data source to query SDRS protection groups within HuaweiCloud.
---

# huaweicloud_sdrs_protection_groups

Use this data source to query SDRS protection groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_protection_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the protection group status.
  For details, see [Protection group status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152930.html).

* `name` - (Optional, String) Specifies the name of a protection group. Fuzzy search is supported.

* `query_type` - (Optional, String) Specifies the query type. Valid values are:
  + **status_abnormal**: Indicates to query protection groups in the abnormal status.
  + **stop_protected**: Indicates to query protection groups for which the protection is disabled.
  + **period_no_dr_drill**: Indicates to query the protection groups for which the no DR drills have been performed in a
    specified duration. The default duration is `3` months.

  This parameter is invalid when the value is set to **general** or left empty.

* `availability_zone` - (Optional, String) Specifies the current production site AZ of a protection group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `server_groups` - The information about protection groups.

  The [server_groups](#server_groups_struct) structure is documented below.

<a name="server_groups_struct"></a>
The `server_groups` block supports:

* `progress` - The synchronization progress of a protection group.

* `source_availability_zone` - The production site AZ configured when a protection group is created.
  The value does not change after a planned failover or failover.

* `domain_id` - The ID of an active-active domain.

* `priority_station` - The current production site of a protection group. Valid values are:
  + **source**: Indicates that the current production site AZ is the `source_availability_zone` value.
  + **target**: Indicates that the current production site AZ is the `target_availability_zone` value.

* `disaster_recovery_drill_num` - The number of DR drills in a protection group.

* `id` - The ID of a protection group.

* `name` - The name of a protection group.

* `description` - The description of a protection group.

* `target_vpc_id` - The ID of the VPC for the DR site.

* `protected_status` - The protected status. Valid values are **started** and **stopped**.

  -> The system has been upgraded. For newly protection groups, the value of this parameter is **null**.

* `protection_type` - The protection mode. Valid values are:
  + **replication-pair**: Indicates that data synchronization is performed at the replication pair level.
  + **null**: Indicates that data synchronization is performed at the replication consistency group level.

  -> The system has been upgraded. Data synchronization is performed at the replication pair level for all resources,
  and the returned value is **replication-pair**.

* `status` - The status of a protection group.
  For details, see [Protection Group Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152930.html)

* `domain_name` - The name of an active-active domain.

* `protected_instance_num` - The number of protected instances in a protection group.

* `health_status` - The health status of a protection group. Valid values are:
  + **normal**: The protection group is normal.
  + **abnormal**: The protection group is abnormal.

  -> The system is upgraded recently. For protection groups created after the upgrade, the value of this parameter is **null**.

* `source_vpc_id` - The ID of the VPC for the production site.

* `created_at` - The time when a protection group was created. The default format is as follows: **yyyy-MM-dd HH:mm:ss.SSS**.

* `updated_at` - The time when a protection group was updated. The default format is as follows: **yyyy-MM-dd HH:mm:ss.SSS**.

* `target_availability_zone` - The DR site AZ configured when a protection group is created.
  The value does not change after a planned failover or failover.

* `replication_num` - The number of replication pairs in a protection group.

* `replication_status` - The data synchronization status. Valid values are:
  + **active**: Data has been synchronized.
  + **inactive**: Data is not synchronized.
  + **copying**: Data is being synchronized.
  + **active-stopped**: Data synchronization is stopped.

  -> The system has been upgraded. For newly protection groups, the value of this parameter is **null**.

* `dr_type` - The deployment model. The default value is **migration**, indicating migration within a VPC.

* `server_type` - The type of managed servers. Valid value is **ECS**, indicates that ECSs are managed.
