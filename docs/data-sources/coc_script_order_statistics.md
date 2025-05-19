---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_order_statistics"
description: |-
  Use this data source to query script order status statistics.
---

# huaweicloud_coc_script_order_statistics

Use this data source to query script order status statistics.

## Example Usage

```hcl
variable "execute_uuid" {}

data "huaweicloud_coc_script_order_statistics" "test" {
  execute_uuid = var.execute_uuid
}
```

## Argument Reference

The following arguments are supported:

* `execute_uuid` - (Required, String) Specifies the execution ID of a script order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `execute_statistics` - Indicates the statistical details.

  The [execute_statistics](#data_execute_statistics_struct) structure is documented below.

<a name="data_execute_statistics_struct"></a>
The `execute_statistics` block supports:

* `instance_status` - Indicates the status of the execution instance.

* `instance_count` - Indicates the number of instances executed in this state.

* `batch_indexes` - Indicates a list of batch indexes in this state.
