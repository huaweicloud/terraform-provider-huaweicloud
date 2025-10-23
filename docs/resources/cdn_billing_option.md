---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_billing_option"
description: |-
  Manages a CDN billing option resource within HuaweiCloud.
---

# huaweicloud_cdn_billing_option

Manages a CDN billing option resource within HuaweiCloud.

## Example Usage

```hcl
variable "charge_mode" {}

resource "huaweicloud_cdn_billing_option" "test" {
  charge_mode  = var.charge_mode
  product_type = "base"
  service_area = "mainland_china"
}
```

## Argument Reference

The following arguments are supported:

* `charge_mode` - (Required, String) Specifies the billing option.  
  The valid values are as follows:
  + **flux**: traffic-based billing.
  + **bw**: bandwidth-based billing, only for V2 and higher customers.

  -> 1. If you change to be billed by peak bandwidth, your traffic package will be frozen. When you switch back,
  the traffic package still takes effect if it is within the required duration.<br> 2. The changes will take effect at
  00:00:00 (GMT+08:00) on the night of the day you submitted the changes. You can continue to make changes until that
  time when your most recently submitted change has taken effect.

* `product_type` - (Required, String) Specifies the product mode.  
  Currently, only **base** (basic acceleration) is supported.

* `service_area` - (Required, String) Specifies the service area.  
  Currentlt, only **mainland_china** (Chinese mainland) is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time, in RFC3339 format.

* `effective_time` - The effective time, in RFC3339 format.

* `status` - The status.

* `current_charge_mode` - The billing option of the account.
