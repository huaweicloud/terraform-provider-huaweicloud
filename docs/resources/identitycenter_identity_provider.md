---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_identity_provider"
description: |-
  Manages an Identity Center identity provider resource within HuaweiCloud.
---

# huaweicloud_identitycenter_identity_provider

Manages an Identity Center identity provider resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "idp_certificate" {}
variable "entity_id" {}
variable "login_url" {}

resource "huaweicloud_identitycenter_identity_provider" "test"{
  identity_store_id = var.identity_store_id
  idp_certificate   = var.idp_certificate
  entity_id         = var.entity_id
  login_url         = var.login_url
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `idp_saml_metadata` - (Optional, String, NonUpdatable) Specifies the metadata of the identity provider.

* `idp_certificate` - (Optional, String, NonUpdatable) Specifies the certificate of the identity provider.

* `entity_id` - (Optional, String) Specifies the ID of the identity provider entity.

* `login_url` - (Optional, String) Specifies the login url of the identity provider.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `idp_certificate_ids` - The list of the identity provider certificate.
  The [idp_certificate_ids](#idp_certificate_ids_struct) structure is documented below.

* `want_request_signed` - Whether the identity provider require the signed request.

* `is_enabled` - Whether the identity provider is enabled.

<a name="idp_certificate_ids_struct"></a>
The `idp_certificate_ids` block supports:

* `certificate_id` - The ID of the identity provider certificate.

* `status` - The status of the identity provider.

## Import

The IdentityCenter identity provider can be imported using the `identity_store_id` and `idp_id`
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_identity_provider.test <identity_store_id>/<idp_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `idp_saml_metadata`, `idp_certificate`.
It is generally recommended running `terraform plan` after importing an IdentityCenter identity provider.
You can then decide if changes should be applied to the IdentityCenter identity provider,
or the resource definition should be updated to align with the identity provider.
Also, you can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_identity_provider" "test" {
  ...

  lifecycle {
    ignore_changes = [
      idp_saml_metadata,
      idp_certificate,
    ]
  }
}
```
