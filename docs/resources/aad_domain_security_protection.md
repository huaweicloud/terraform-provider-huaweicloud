---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domain_security_protection"
description: |-
  Use this resource to modify Advanced Anti-DDos security protection within HuaweiCloud.
---

# huaweicloud_aad_domain_security_protection

Use this resource to modify Advanced Anti-DDos security protection within HuaweiCloud.

-> This resource is only a one-time action resource for updating AAD security protection. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "domain_id" {}
variable "waf_switch" {}
variable "cc_switch" {}

resource "huaweicloud_aad_domain_security_protection" "test" {
  domain_id  = var.domain_id
  waf_switch = var.waf_switch
  cc_switch  = var.cc_switch
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String, NonUpdatable) Specifies the domain ID.

* `waf_switch` - (Required, Int, NonUpdatable) Specifies whether to enable basic web protection. Valid values are:
  + `0`: On.
  + `1`: Off.

* `cc_switch` - (Required, Int, NonUpdatable) Specifies whether to enable CC protection. Valid values are:
  + `0`: On.
  + `1`: Off.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `domain_id`).
