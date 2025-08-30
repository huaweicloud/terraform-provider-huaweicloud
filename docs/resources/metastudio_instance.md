---
subcategory: "Meta Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_metastudio_instance"
description: |-
  Manages a Meta Studio instance resource within HuaweiCloud.
---

# huaweicloud_metastudio_instance

Manages a Meta Studio instance resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_metastudio_instance" "resource1_periodic" {
  period_type        = 2
  period_num         = 1
  is_auto_renew      = 0
  resource_spec_code = "hws.resource.type.metastudio.modeling.avatarlive.channel"
}

resource "huaweicloud_metastudio_instance" "resource2_one_time" {
  period_type        = 6
  period_num         = 1
  is_auto_renew      = 0
  resource_spec_code = "hws.resource.type.metastudio.avatarmodeling.number"
}
```

~> Terraform will automatically sign the auto-pay agreement during resource creation.
  No user intervention is required for this process

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `period_type` - (Required, Int, NonUpdatable) Specifies the charging period unit.
  Valid values are:
  + `2` - Month
  + `3` - Year
  + `6` - One Time. When period_type is 6, the period_num should be 1

* `period_num` - (Required, Int, NonUpdatable) Specifies the number of periods to purchase.
  Value range: `1` to `2147483647`.

* `is_auto_renew` - (Optional, Int, NonUpdatable) Specifies whether to auto-renew the resource when it expires.
  Valid values are:
  + `0` - Do not renew automatically (default)
  + `1` - Renew automatically

* `resource_spec_code` - (Required, String, NonUpdatable) Specifies the resource specification code for
  user-purchased cloud service products.
  For details, see [Resource Types](https://support.huaweicloud.com/api-metastudio/metastudio_02_0042.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `order_id` - The order ID associated with the resource.

* `resource_expire_time` - The expiration time of the resource.

* `business_type` - The business type of the resource.

* `sub_resource_type` - The sub-resource type.

* `is_sub_resource` - Indicates whether it is a sub-resource.

* `charging_mode` - The billing mode of the resource (e.g., `PERIODIC` for pre-paid).

* `amount` - The total amount of the resource.

* `usage` - The usage amount of the resource.

* `status` - The status of the resource:
  + `0` - Normal
  + `1` - Frozen

* `unit` - The unit of measurement for the resource amount.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes (includes order processing time).
* `delete` - Default is 10 minutes.
