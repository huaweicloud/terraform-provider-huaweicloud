---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_order_batches"
description: |-
  Use this data source to query the batch list of script orders.
---

# huaweicloud_coc_script_order_batches

Use this data source to query the batch list of script orders.

## Example Usage

```hcl
variable "execute_uuid" {}

data "huaweicloud_coc_script_order_batches" "test" {
  execute_uuid = var.execute_uuid
}
```

## Argument Reference

The following arguments are supported:

* `execute_uuid` - (Required, String) Specifies the execution ID of a script order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the batch list of script orders.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `batch_index` - Indicates the batch index.

* `total_instances` - Indicates the number of instance nodes in the batch.

* `rotation_strategy` - Indicates suspension and resumption policy.

* `properties` - Indicates the batch label.
