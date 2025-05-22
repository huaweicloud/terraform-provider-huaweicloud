---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_web_protection_policies"
description: |-
  Use this data source to get the list of Advanced Anti-DDos web protection policies within HuaweiCloud.
---

# huaweicloud_aad_web_protection_policies

Use this data source to get the list of Advanced Anti-DDos web protection policies within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "overseas_type" {}

data "huaweicloud_aad_web_protection_policies" "test" {
  domain_name   = var.domain_name
  overseas_type = var.overseas_type
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name.

* `overseas_type` - (Required, Int) Specifies the overseas type. Valid values are `0` (Mainland) and `1` (Overseas).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `options` - The protection options.

  The [options](#options_struct) structure is documented below.

* `level` - The intelligent CC protection level. Valid values are:
  + `0`: Loose.
  + `1`: Normal.
  + `2`: Strict.

* `mode` - The smart CC mode. Valid values are:
  + `0`: Early warning.
  + `1`: Protection.

<a name="options_struct"></a>
The `options` block supports:

* `geoip` - Whether to enable regional ban protection.

* `whiteblackip` - Whether to enable blacklist and whitelist protection.

* `modulex_enabled` - Whether to enable intelligent CC protection.

* `cc` - Whether to enable CC (frequency control).

* `custom` - Whether to enable precise access protection.
