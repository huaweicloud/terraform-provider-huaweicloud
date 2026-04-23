---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_execution_plan_prices"
description: |-
  Use this dataSource to get the price estimate of an execution plan.
---

# huaweicloud_rfs_execution_plan_prices

Use this dataSource to get the price estimate of an execution plan.

## Example Usage

```hcl
variable "stack_name" {}
variable "execution_plan_name" {}

data "huaweicloud_rfs_execution_plan_prices" "test" {
  stack_name          = var.stack_name
  execution_plan_name = var.execution_plan_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `stack_name` - (Required, String) Specifies the name of the resource stack.

* `execution_plan_name` - (Required, String) Specifies the name of the execution plan.

* `stack_id` - (Optional, String) Specifies the ID of the resource stack.

* `execution_plan_id` - (Optional, String) Specifies the ID of the execution plan.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `currency` - The currency. For example, **CNY** is returned for the Chinese site.

* `items` - The price estimate of all resources in the execution plan.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `resource_type` - The resource type in the template.

* `resource_name` - The logical resource name in the template, used as the default resource name.

* `index` - The index of the resource. If `count` or `for_each` is used, this value is set.

* `module_address` - The module address of the resource.

* `supported` - Whether the resource or the given parameters support price inquiry.

* `unsupported_message` - The reason the resource does not support price inquiry.

* `resource_price` - The price information for the resource. If the resource is free, this block may be absent.

  The [resource_price](#resource_price_struct) structure is documented below.

<a name="resource_price_struct"></a>
The `resource_price` block supports:

* `charge_mode` - The billing mode.  
  The valid values are **PRE_PAID** (Yearly/Monthly), **POST_PAID** (Pay-per-use), and **FREE** (Free).

* `sale_price` - The final amount after all discounts, excluding promotion discounts and coupons, in **CNY**.

* `discount` - The total discount amount, in **CNY**.

* `original_price` - The original price, in **CNY**.

* `period_type` - The billing period unit, used together with `period_count`.  
  The valid values include **HOUR**, **DAY**, **MONTH**, **YEAR**, **BYTE**, **MB**, and **GB**.

* `period_count` - The number of billing periods.

  For a detailed explanation of this parameter,
  please refer to [documentation](https://support.huaweicloud.com/intl/en-us/api-aos/EstimateExecutionPlanPrice.html).

* `best_discount_type` - The best discount type.

  For a detailed explanation of this parameter,
  please refer to [documentation](https://support.huaweicloud.com/intl/en-us/api-aos/EstimateExecutionPlanPrice.html).

* `best_discount_price` - The best discount amount, in **CNY**.

* `official_website_discount_price` - The official website discount amount, in **CNY**.
