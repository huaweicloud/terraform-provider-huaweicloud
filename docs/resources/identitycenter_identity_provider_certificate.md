---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_identity_provider_certificate"
description: |-
  Manages an Identity Center identity provider certificate resource within HuaweiCloud.
---

# huaweicloud_identitycenter_identity_provider_certificate

Manages an Identity Center identity provider certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "idp_id" {}
variable "x509_certificate_in_pem" {}

resource "huaweicloud_identitycenter_identity_provider_certificate" "test"{
  identity_store_id       = var.identity_store_id
  idp_id                  = var.idp_id
  x509_certificate_in_pem = var.x509_certificate_in_pem
  certificate_use         = "SIGNING"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `idp_id` - (Required, String, NonUpdatable) Specifies the ID of the identity provider.

* `x509_certificate_in_pem` - (Required, String, NonUpdatable) Specifies the certificate of the identity provider.

* `certificate_use` - (Required, String, NonUpdatable) Specifies how the identity provider certificate is used.
Value options: **SIGNING**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `issuer_name` - The issuer of the identity provider certificate.

* `not_after` - The validity period of the identity provider certificate.

* `not_before` - The validity period of the identity provider certificate.

* `public_key` - The public key of the identity provider certificate.

* `serial_number_string` - The serial number of the identity provider certificate.

* `signature_algorithm_name` - The signature algorithm of the identity provider certificate.

* `subject_name` - The subject of the identity provider certificate.

* `version` - The version of the identity provider certificate.

## Import

The IdentityCenter identity provider certificate can be imported using
the `identity_store_id`, `idp_id` and `certificate_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_identity_provider_certificate.test <identity_store_id>/<idp_id>/<certificate_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `certificate_use`. It is generally
recommended running `terraform plan` after importing an IdentityCenter identity provider certificate.
You can then decide if changes should be applied to the IdentityCenter identity provider certificate,
or the resource definition should be updated to align with the identity provider certificate. Also,
you can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_identity_provider_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      certificate_use,
    ]
  }
}
```
