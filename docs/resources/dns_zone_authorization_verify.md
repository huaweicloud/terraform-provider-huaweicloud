---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_authorization_verify"
description: |-
  Using this resource to verify a sub-domain authorization within HuaweiCloud.
---

# huaweicloud_dns_zone_authorization_verify

Using this resource to verify a sub-domain authorization within HuaweiCloud.

-> This resource is only a one-time action resource for verify a sub-domain after the authorization request is send.
   Deleting this resource will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

```hcl
variable "sub_domain_authorization_id" {}

resource "huaweicloud_dns_zone_authorization_verify" "test" {
  authorization_id = var.sub_domain_authorization_id
}
```

## Argument Reference

The following arguments are supported:

* `authorization_id` - (Required, String, NonUpdatable) Specifies the ID of the sub-domain authorization to be verified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The authorization status for this verification.
  + **CREATED**: Authorization has been created.
  + **VERIFIED**: Authorization has been verified.
