---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_orders"
description: |-
  Use this data source to query the list of script orders.
---

# huaweicloud_coc_script_orders

Use this data source to query the list of script orders.

## Example Usage

```hcl
data "huaweicloud_coc_script_orders" "test" {}
```

## Argument Reference

The following arguments are supported:

* `start_time` - (Optional, Int) Specifies the start time.

* `end_time` - (Optional, Int) Specifies the end time.

* `creator` - (Optional, String) Specifies the creator.

* `status` - (Optional, String) Specifies the script order status.
  Values can be as follows:
  + **READY**: Prepare.
  + **PROCESSING**: The operation is in progress.
  + **ABNORMAL**: Abnormal.
  + **PAUSED**: Paused.
  + **CANCELED**: Canceled.
  + **FINISHED**: Success.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of script orders.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `order_id` - Indicates the primary key ID.

* `order_name` - Indicates the script order name.

* `execute_uuid` - Indicates this UUID is used when the list is redirected to the details page.

* `gmt_created` - Indicates the creation time.

* `gmt_finished` - Indicates the completion time.

* `execute_costs` - Indicates the execution duration in seconds.

* `creator` - Indicates the creator.

* `status` - Indicates the script order status.

* `properties` - Indicates the label.

  The [properties](#data_properties_struct) structure is documented below.

<a name="data_properties_struct"></a>
The `properties` block supports:

* `region_ids` - Indicates the region ID of the Cloud CMDB service instance.
