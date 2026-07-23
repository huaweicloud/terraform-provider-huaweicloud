---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_alarm_overview"
description: |-
  Use this data source to get the alarm data statistics overview within HuaweiCloud.
---

# huaweicloud_dsc_alarm_overview

Use this data source to get the alarm data statistics overview within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_alarm_overview" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_source_info` - The statistics quantity of each alarm source.

  The [alarm_source_info](#alarm_source_info_struct) structure is documented below.

* `total_alarm` - The statistics of all alarms grouped by level.

  The [total_alarm](#total_alarm_struct) structure is documented below.

* `turn_off_num` - The total number of alarms in the closed status.

* `turn_on_num` - The total number of alarms in the enabled status.

* `untreated_alarm` - The statistics of untreated alarms grouped by level.

  The [untreated_alarm](#untreated_alarm_struct) structure is documented below.

<a name="alarm_source_info_struct"></a>
The `alarm_source_info` block supports:

* `api_num` - The number of alarms from API instances.

* `cbh_num` - The number of alarms from CBH instances.

* `database_encrypt_num` - The number of alarms from database encryption instances.

* `database_op_num` - The number of alarms from database operation instances.

* `dbss_num` - The number of alarms from database audit (DBSS) instances.

<a name="total_alarm_struct"></a>
The `total_alarm` block supports:

* `fatal_num` - The number of fatal alarms.

* `high_risk_num` - The number of high risk alarms.

* `middle_risk_num` - The number of middle risk alarms.

* `low_risk_num` - The number of low risk alarms.

* `notice_num` - The number of notice alarms.

<a name="untreated_alarm_struct"></a>
The `untreated_alarm` block supports:

* `fatal_num` - The number of fatal alarms.

* `high_risk_num` - The number of high risk alarms.

* `middle_risk_num` - The number of middle risk alarms.

* `low_risk_num` - The number of low risk alarms.

* `notice_num` - The number of notice alarms.
