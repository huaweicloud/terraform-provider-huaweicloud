---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_billing_option"
description: |-
  Use this data source to get CDN billing option within HuaweiCloud.
---

# huaweicloud_cdn_billing_option

Use this data source to get CDN billing option within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_billing_option" "test" {
  product_type = "base"
}
```

## Argument Reference

The following arguments are supported:

* `product_type` - (Required, String) Specifies the product mode.  
  Currently, only **base** (basic acceleration) is supported.

* `service_area` - (Optional, String) Specifies the service area.  
  The valid values are as follows:
  + **mainland_china**: Chinese mainland.
  + **outside_mainland_china**: outside the Chinese mainland.

  Defaults to **mainland_china**.

* `status` - (Optional, String) Specifies the billing option status.  
  The valid values are as follows:
  + **active**: effective.
  + **upcoming**: to take effect.

  Defaults to **active**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `charge_mode` - The billing option.  
  + **flux**: traffic-based billing.
  + **bw**: bandwidth-based billing, only for V2 and higher customers.

* `created_at` - The creation time, in RFC3339 format.

* `effective_time` - The effective time, in RFC3339 format.
