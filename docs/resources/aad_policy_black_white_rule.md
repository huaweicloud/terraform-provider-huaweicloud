---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_policy_black_white_rule"
description: |-
  Manages a WAF black and white rule resource within HuaweiCloud AAD service.
---

# huaweicloud_aad_policy_black_white_rule

Manages a WAF black and white rule resource within HuaweiCloud AAD service.

## Example Usage

```hcl
variable "domain_name" {}
variable "ip" {}

resource "huaweicloud_aad_policy_black_white_rule" "white_rule" {
  domain_name   = var.domain_name
  ip            = var.ip
  overseas_type = 0
  type          = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, NonUpdatable) Specifies the domain name.

* `ip` - (Required, String, NonUpdatable) Specifies the IP address or IP segment.

* `overseas_type` - (Required, Int, NonUpdatable) Specifies the protection area. The value can be:
  + `0`: Chinese mainland.
  + `1`: Outside Chinese mainland.

* `type` - (Required, Int, NonUpdatable) Specifies the rule type. The value can be:
  + `0`: Blacklist.
  + `1`: Whitelist.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (rule ID).

* `domain_id` - The domain ID.

## Import

The AAD policy black white rule can be imported using the `domain_name`, `overseas_type`, `ip` and `type`,
separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_aad_policy_black_white_rule.test <domain_name>/<overseas_type>/<ip>/<type>
```
