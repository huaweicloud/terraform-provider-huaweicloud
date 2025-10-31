---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_bearer_token"
description: |-
  Manages an Identity Center bearer token resource within HuaweiCloud.
---

# huaweicloud_identitycenter_bearer_token

Manages an Identity Center bearer token resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "tenant_id" {}

resource "huaweicloud_identitycenter_bearer_token" "test"{
  identity_store_id = var.identity_store_id
  tenant_id         = var.tenant_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `tenant_id` - (Required, String, NonUpdatable) Specifies the ID of the tenant.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creation_time` - The creation time of the bearer token.

* `expiration_time` - The expiration time of the bearer token.

## Import

The IdentityCenter bearer token can be imported using the `identity_store_id` `tenant_id` and `token_id`
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_bearer_token.test <identity_store_id>/<tenant_id>/<token_id>
```
