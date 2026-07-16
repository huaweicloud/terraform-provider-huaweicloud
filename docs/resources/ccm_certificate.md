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

### Single Domain Free Certificate

```hcl
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand     = "SYMANTIC"
  cert_type      = "DV_SSL_CERT_BASIC"
  domain_type    = "SINGLE_DOMAIN"
  effective_time = 0
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

* `cert_brand` - (Required, String, ForceNew) Specifies the certificate authority.
  The valid values are as follows:
  + **GEOTRUST** - Use when you need a cost-effective certificate with good compatibility.  
    Supports DV, OV, and EV types. DV and OV support single-domain, multi-domain, and wildcard;  
    EV supports single-domain and multi-domain. OV also supports IP address binding.
  + **GLOBALSIGN** - Use when you need fast certificate issuance and low server resource consumption.  
    Only supports OV and EV types; does not support DV. OV supports single-domain, multi-domain,
    wildcard, and IP address binding;  
    EV supports single-domain and multi-domain.
  + **SYMANTEC** - Use when you need the highest industry trust level, especially for banking,
    finance, and high-security digital transaction scenarios.  
    Supports DV, OV, and EV types, including Pro editions.  
    OV and EV Pro support single-domain, multi-domain, and wildcard;  
    EV supports single-domain and multi-domain.
  + **CFCA** - Use when you need a domestic (Chinese) certificate with SM2 national cryptographic
    algorithm support, or when financial-grade security services are required.  
    Only supports OV and EV types; does not support DV. OV supports single-domain, multi-domain, and wildcard;  
    EV supports single-domain and multi-domain.
  + **TRUSTASIA** - Use when you need a domestic (Chinese) certificate with SM2 national cryptographic
    algorithm support.  
    Only supports OV and EV types; does not support DV.
  + **VTRUS** - Use when you need a domestic (Chinese) certificate with SM2 national cryptographic
    algorithm support and need compatibility with national cryptographic browsers.  
    Supports DV and OV types; does not support EV. DV supports single-domain and wildcard;  
    OV supports single-domain, multi-domain, wildcard, and IP address binding.

  Changing this parameter will create a new resource.

* `cert_type` - (Required, String, ForceNew) Specifies the certificate type.
  The valid values are as follows:
  + **DV_SSL_CERT** - Domain Validated certificate.  
    Use for personal websites, development testing, or any scenario where only domain ownership
    needs to be verified. Security level is general; issuance takes hours.  
    Only verifies domain ownership, no organizational identity validation.
  + **DV_SSL_CERT_BASIC** - Domain Validated Basic certificate.  
    Use for simple personal websites or development/testing environments with minimal requirements.  
    This is the most basic certificate option with limited domain type support (single-domain only).
  + **EV_SSL_CERT** - Extended Validation certificate.  
    Use for large enterprises and organizations with strict security requirements, such as finance, insurance,
    and banking.  
    Highest security level with strict organizational identity validation; issuance takes `7`-`10` business days.  
    Supports single-domain and multi-domain.
  + **EV_SSL_CERT_PRO** - Extended Validation Pro certificate.  
    Use when you need the highest security level with premium features such as additional domain types
    or enhanced encryption algorithm support (including ECC).  
    Same strict organizational validation as EV, with more domain type and encryption options.
  + **OV_SSL_CERT** - Organization Validated certificate.  
    Use for education, government, Internet, SME, and e-commerce websites (e.g., Apple Store, WeChat Mini Programs).  
    High security level with comprehensive organizational identity validation; issuance takes `3`-`5` business days.  
    Supports single-domain, multi-domain, wildcard, and IP address binding.
  + **OV_SSL_CERT_PRO** - Organization Validated Pro certificate.  
    Use when you need high security with premium features such as enhanced encryption algorithm
    support (including ECC) or additional domain type options.  
    Same organizational validation as OV, with more encryption and domain configuration options.

  Changing this parameter will create a new resource.

* `domain_type` - (Required, String, ForceNew) Specifies the type of domain name.  
  The valid values are as follows:
  + **SINGLE_DOMAIN**
  + **MULTI_DOMAIN**
  + **WILDCARD**

  Changing this parameter will create a new resource.

* `effective_time` - (Required, Int, ForceNew) Specifies the validity period (year).  
  The valid values are `0`, `1`, `2`, and `3`.  
  Changing this parameter will create a new resource.

  -> The value `0` can only be used when applying for a free CCM certificate, and its validity period is 3 months.

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

* `create` - Default is `10` minutes.
* `update` - Default is `10` minutes.
* `delete` - Default is `10` minutes.

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
