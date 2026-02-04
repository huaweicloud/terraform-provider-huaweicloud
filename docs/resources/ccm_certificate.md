---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate"
description: |-
  Manages a CCM SSL certificate resource within HuaweiCloud.
---

# huaweicloud_ccm_certificate

Manages a CCM SSL certificate resource within HuaweiCloud.

-> Refer to [document](https://support.huaweicloud.com/intl/en-us/productdesc-ccm/ccm_01_0219.html) to see the
differences between different types of certificates.

## Example Usage

### Single Domain Certificate

```hcl
variable "cert_brand" {}
variable "cert_type" {}

resource "huaweicloud_ccm_certificate" "test" {
  cert_brand     = var.cert_brand
  cert_type      = var.cert_type
  domain_type    = "SINGLE_DOMAIN"
  effective_time = 1
  domain_numbers = 1
}
```

### Multi Domain Certificate

```hcl
variable "cert_brand" {}
variable "cert_type" {}

resource "huaweicloud_ccm_certificate" "test" {
  cert_brand             = var.cert_brand
  cert_type              = var.cert_type
  domain_type            = "MULTI_DOMAIN"
  effective_time         = 1
  domain_numbers         = 4
  primary_domain_type    = "SINGLE_DOMAIN"
  single_domain_number   = 1
  wildcard_domain_number = 2
}
```

### Wildcard Domain Certificate

```hcl
variable "cert_brand" {}
variable "cert_type" {}

resource "huaweicloud_ccm_certificate" "test" {
  cert_brand     = var.cert_brand
  cert_type      = var.cert_type
  domain_type    = "WILDCARD"
  effective_time = 1
  domain_numbers = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cert_brand` - (Required, String, ForceNew) Specifies the certificate authority. Valid values are **GEOTRUST**,
  **GLOBALSIGN**, **SYMANTEC**, **CFCA**, **TRUSTASIA**, and **VTRUS**.

  Changing this parameter will create a new resource.

* `cert_type` - (Required, String, ForceNew) Specifies the certificate type. Valid values are **DV_SSL_CERT**,
  **DV_SSL_CERT_BASIC**, **EV_SSL_CERT**, **EV_SSL_CERT_PRO**, **OV_SSL_CERT**, and **OV_SSL_CERT_PRO**.

  Changing this parameter will create a new resource.

* `domain_type` - (Required, String, ForceNew) Specifies the type of domain name. Valid values are **SINGLE_DOMAIN**,
  **MULTI_DOMAIN**, and **WILDCARD**.

  Changing this parameter will create a new resource.

* `effective_time` - (Required, Int, ForceNew) Specifies the validity period (year). Valid values are `1`, `2`, and `3`.

  Changing this parameter will create a new resource.

* `domain_numbers` - (Required, Int, ForceNew) Specifies the quantity of domain name.
  + When `domain_type` is set to **SINGLE_DOMAIN** or **WILDCARD**, this field can only be set to `1`.
  + When `domain_type` is set to **MULTI_DOMAIN**, the value of this field ranges from `2` to `250`. The value of this
    field should be the number of additional domain names plus one main domain name.
    For example, if field `single_domain_number` is set to `1`, and field `wildcard_domain_number` is set to `2`, then
    the value of this field should be `4`.

  Changing this parameter will create a new resource.

* `primary_domain_type` - (Optional, String, ForceNew) Specifies the type of primary domain name in multiple domains.
  Valid values are **SINGLE_DOMAIN** and **WILDCARD_DOMAIN**.

  Changing this parameter will create a new resource.

* `single_domain_number` - (Optional, Int, ForceNew) Specifies the number of additional single domain names.
  The value of this field ranges from `1` to `249`.

  Changing this parameter will create a new resource.

* `wildcard_domain_number` - (Optional, Int, ForceNew) Specifies the number of additional wildcard domain names.
  The value of this field ranges from `0` to `248`.

  Changing this parameter will create a new resource.

-> Fields `primary_domain_type`, `single_domain_number`, and `wildcard_domain_number` are required when `domain_type`
is set to **MULTI_DOMAIN**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the SSL certificate ID.

* `validity_period` - The validity period (month).

* `status` - The certificate status. Valid values are:
  + **PAID**: The certificate has been paid and the certificate is waiting to be applied for.
  + **ISSUED**: The certificate has been issued.
  + **CHECKING**: The certificate application is under review.
  + **CANCELCHECKING**: The certificate cancellation application is under review.
  + **UNPASSED**: The certificate application failed.
  + **EXPIRED**: The certificate has expired.
  + **REVOKING**: The certificate revocation application is under review.
  + **REVOKED**: The certificate has been revoked.
  + **UPLOAD**: The certificate in custody.
  + **SUPPLEMENTCHECKING**: Additional domain names added to multi-domain certificates are under review.
  + **CANCELSUPPLEMENTING**: Cancel the addition of additional domain names under review.

* `order_id` - The order ID.

* `name` - The certificate name.

* `push_support` - Whether the certificate supports push.

* `revoke_reason` - The reason for certificate revocation.

* `signature_algorithm` - The signature algorithm.

* `issue_time` - The certificate issuance time.

* `not_before` - The certificate validity time.

* `not_after` - The certificate expiration time.

* `validation_method` - The authentication method of domain name.

* `domain` - The domain name bound to the certificate.

* `sans` - The information of additional domain name for the bound certificate.

* `fingerprint` - The SHA-1 fingerprint of the certificate.

* `authentification` - The ownership certification information of domain name.
  The [authentification](#CCMCertificate_authentification) structure is documented below.

<a name="CCMCertificate_authentification"></a>
The `authentification` block supports:

* `record_name` - The name of the domain name check value.

* `record_type` - The type of the domain name check value.

* `record_value` - The domain name check value.

* `domain` - The domain name corresponding to the check value.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The CCM certificate can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ccm_certificate.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `effective_time`, `single_domain_number`,
`tags`. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_ccm_certificate" "test" { 
  ...

  lifecycle {
    ignore_changes = [
      effective_time, single_domain_number,
    ]
  }
}
```
