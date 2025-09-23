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

* `product_type` - (Required, String) Specifies the product mode. Only **base** (basic acceleration) is supported.

* `status` - (Optional, String) Specifies the billing option status. Valid values are **active** (effective) and
  **upcoming** (to take effect). Defaults to **active**.

* `service_area` - (Optional, String) Specifies the service area. Valid values are **mainland_china** (Chinese mainland)
  and **outside_mainland_china** (outside the Chinese mainland). Defaults to **mainland_china**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `charge_mode` - Indicates the billing option. The value can be **flux** (traffic-based billing) or
  **bw** (bandwidth-based billing, only for V2 and higher customers).

* `created_at` - Indicates the creation time.

* `effective_time` - Indicates the effective time of the option.
