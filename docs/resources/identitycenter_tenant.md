---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_tenant"
description: |-
  Manages an Identity Center tenant resource within HuaweiCloud.
---

# huaweicloud_identitycenter_tenant

Manages an Identity Center tenant resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}

resource "huaweicloud_identitycenter_tenant" "test"{
  identity_store_id = var.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creation_time` - The creation name of the tenant.

* `scim_endpoint` - The SCIM endpoint.

## Import

The IdentityCenter tenant can be imported using the `identity_store_id` and `tenant_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_tenant.test <identity_store_id>/<tenant_id>
```
