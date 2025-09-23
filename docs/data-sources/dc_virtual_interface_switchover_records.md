---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_interface_switchover_records"
description: |-
  Use this data source to get the list of DC virtual interface switchover test records.
---

# huaweicloud_dc_virtual_interface_switchover_records

Use this data source to get the list of DC virtual interface switchover test records.

## Example Usage

```hcl
data "huaweicloud_dc_virtual_interface_switchover_records" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Optional, List) Specifies the resource ID used for querying switchover test records.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `switchover_test_records` - Indicates the list of the switchover test records.

  The [switchover_test_records](#switchover_test_records_struct) structure is documented below.

<a name="switchover_test_records_struct"></a>
The `switchover_test_records` block supports:

* `id` - Indicates the unique ID of the switchover test record.

* `resource_id` - Indicates the ID of the resource on which the switchover test is to be performed.

* `resource_type` - Indicates the type of the resource on which the switchover test is to be performed.

* `operation` - Indicates whether to perform a switchover test.
  The value can be: **shutdown** and **undo_shutdown**.

* `operate_status` - Indicates the switchover test status.
  The value can be:
  + **STARTING**: indicates the initial status.
  + **INPROGRESS**: The configuration is being delivered.
  + **COMPLETE**: The configuration is delivered.
  + **ERROR**: The configuration fails to be delivered.

* `start_time` - Indicates the start time of the switchover test.

* `end_time` - Indicates the end time of the switchover test.
