---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_asymmetric_signature_switch"
description: |-
  Use this resource to enable or disable the asymmetric signature function for the account within HuaweiCloud.
---

# huaweicloud_identityv5_asymmetric_signature_switch

Use this resource to enable or disable the asymmetric signature function for the account within HuaweiCloud.

-> This resource is a one-time action resource used to enable or disable the asymmetric signature function. Deleting
   this resource will not clear the corresponding request record, but will only remove the resource information from
   the tf state file.

## Example Usage

```hcl
resource "huaweicloud_identityv5_asymmetric_signature_switch" "test" {
  asymmetric_signature_switch = true
}
```

## Argument Reference

* `asymmetric_signature_switch` - (Required, Bool) Specifies Whether to enable the asymmetric signature function.

## Attribute Reference

* `id` - The ID of the resource, which is the domain ID.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_asymmetric_signature_switch.test <id>
```
