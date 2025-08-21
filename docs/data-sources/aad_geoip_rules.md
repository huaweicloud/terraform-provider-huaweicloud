---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_geoip_rules"
description: |-
  Use this data source to get the list of Geo IP rules within HuaweiCloud.
---

# huaweicloud_aad_geoip_rules

Use this data source to get the list of Geo IP rules within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_aad_geoip_rules" "test" {
  domain_name   = var.domain_name
  overseas_type = "0"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name to query.

* `overseas_type` - (Required, String) Specifies the protection region.
  + `0`: Mainland China.
  + `1`: Outside Mainland China.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - The list of Geo IP rules.
  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `name` - The name of the rule.

* `geoip` - The geographical location code.

* `overseas_type` - The protection region.

* `timestamp` - The creation timestamp of the rule.

* `white` - The protection action. The options are as follows:
  + `0`: Block.
  + `1`: Allow.
  + `2`: Log only.
