---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_asymmetric_signature_switch"
description: |-
  Enable or disable the asymmetric signature function for the account within HuaweiCloud.
---

# huaweicloud_identityv5_asymmetric_signature_switch

Enable or disable the asymmetric signature function for the account within HuaweiCloud.

->**Note** The asymmetric signature switch can not be destroyed.

## Example Usage

```hcl
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = true
}
```

## Argument Reference

* `asymmetric_signature_switch` - (Required, Bool) Specifies the asymmetric signature switch.

## Attribute Reference

* `id` - Resource ID in format `<domain_id>`.
