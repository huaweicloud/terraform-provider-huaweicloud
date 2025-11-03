---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_user_quotas"
description: |-
  Use this data source to get the Advanced Anti-DDos user quotas within HuaweiCloud.
---

# huaweicloud_aad_user_quotas

Use this data source to get the Advanced Anti-DDos user quotas within HuaweiCloud.

## Example Usage

```hcl
variable "type" {}

data "huaweicloud_aad_user_quotas" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String) Specifies the request type.
  The valid values are **instance**, **domain**, **port**, **waf**, and **domain_port**.

* `overseas_type` - (Optional, String) Specifies the protection region.
  This parameter is mandatory when `type` is set to **domain**.

* `ip` - (Optional, String) Specifies the high-defense IP. This parameter is mandatory when `type` is set to **port**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `domain` - The remaining domain quota.

* `instance` - The remaining instance quota.

* `port` - The remaining forwarding configuration quota for the specified IP.

* `domain_port_quota` - The domain server configuration quota.

* `cc_quota` - The WAF CC rule quota.

* `custom` - The WAF precise protection rule quota.

* `geo_ip` - The WAF regional protection quota.

* `white_ip` - The WAF whitelist quota.
