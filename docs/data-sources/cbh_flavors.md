---
subcategory: "Cloud Bastion Host (CBH)"
---

# huaweicloud_cbh_flavors

Use this data source to get the list of CBH flavor.

## Example Usage

```hcl
variable "project_id" {}
variable "flavor_1" {
  type = object({
    resource_spec    = string
    region           = string
    period_unit      = string
    period           = number
    subscription_num = number
  })
}
variable "flavor_2" {
  type = object({
    resource_spec    = string
    region           = string
    period_unit      = string
    period           = number
    subscription_num = number
  })
}

data "huaweicloud_cbh_flavors" "test" {
  project_id = var.project_id

  flavors {
    resource_spec    = var.flavor_1.resource_spec
    region           = var.flavor_1.region
    period_unit      = var.flavor_1.period_unit
    period           = var.flavor_1.period
    subscription_num = var.flavor_1.subscription_num
  }

  flavors {
    resource_spec    = var.flavor_2.resource_spec
    region           = var.flavor_2.region
    period_unit      = var.flavor_2.period_unit
    period           = var.flavor_2.period
    subscription_num = var.flavor_2.subscription_num
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the ID of the project.

* `flavors` - (Required, List) Indicates the list of the product info.
  The [Flavor](#CbhFlavors_Flavor) structure is documented below.

<a name="CbhFlavors_Flavor"></a>
The `Flavor` block supports:

* `resource_spec` - (Required, String) Specifies the resource specifications of cloud service types.

* `region` - (Required, String) Specifies the cloud service region code.

* `period_unit` - (Required, String) Specifies the period type of yearly/monthly product order.
  Value options: **month**, **year**.

* `period` - (Required, Int) Specifies the number of periods of a yearly/monthly product order.

* `subscription_num` - (Required, Int) Specifies the number of subscriptions of a yearly/monthly product order.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `flavors` - Indicates the list of the product info.
  The [Flavor](#CbhFlavors_Flavor) structure is documented below.

<a name="CbhFlavors_Flavor"></a>
The `Flavor` block supports:

* `flavor_id` - Indicates the ID of the product.
