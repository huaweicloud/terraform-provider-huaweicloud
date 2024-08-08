---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_apply"
description: |-
  Manages a CCM SSL certificate apply resource within HuaweiCloud.
---

# huaweicloud_ccm_certificate_apply

Manages a CCM SSL certificate apply resource within HuaweiCloud.

-> 1. The application for a certificate needs to wait for approval before it can be used.
<br/>2. The certificate application results can be obtained through `status` in datasource `huaweicloud_ccm_certificates`.
<br/>3. The current resource is a one-time resource, and destroying this resource will not affect the result of the
certificate application.

## Example Usage

```hcl
variable "certificate_id" {}
variable "domain" {}
variable "applicant_name" {}
variable "applicant_phone" {}
variable "applicant_email" {}
variable "domain_method" {}

resource "huaweicloud_ccm_certificate_apply" "test" {
  certificate_id  = var.certificate_id
  domain          = var.domain
  applicant_name  = var.applicant_name
  applicant_phone = var.applicant_phone
  applicant_email = var.applicant_email
  domain_method   = var.domain_method
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `certificate_id` - (Required, String, ForceNew) Specifies the CCM SSL certificate ID.

  Changing this parameter will create a new resource.

* `domain` - (Required, String, ForceNew) Specifies the domain name bound to the certificate.
  + When the certificate is **single domain name** or **wide domain name**, just fill in the domain name directly.
  + When the certificate is **multiple domain names**, the written domain name will exist as the main domain name.

  Changing this parameter will create a new resource.

* `applicant_name` - (Required, String, ForceNew) Specifies the name of the applicant. The valid length is limited
  from `4` to `100`, only chinese and english letters, underscores (_), hyphens (-), commas (,) and dots (.) are allowed.

  Changing this parameter will create a new resource.

* `applicant_phone` - (Required, String, ForceNew) Specifies the phone number of the applicant.

  Changing this parameter will create a new resource.

* `applicant_email` - (Required, String, ForceNew) Specifies the email of the applicant.

  Changing this parameter will create a new resource.

* `domain_method` - (Required, String, ForceNew) Specifies the domain name verification method. Valid values are:
  + **DNS**: Verify using DNS. Refers to verifying domain name ownership by parsing specified DNS records on the domain
    name management platform.
  + **FILE**: Verify using file. Refers to verifying domain name ownership by creating a specified file on the server.
  + **EMAIL**: Verify using email. Refers to logging in to the domain name administrator's email, receiving the domain
    name confirmation email and following the prompts to verify domain name ownership.

  Changing this parameter will create a new resource.

  -> 1. DV and DV (Basic) certificates are verified by **DNS** by default.
  <br/>2. Pure IP (public IP) certificates only support verification through **FILE**.
  <br/>3. Only pure IP certificates support **FILE** verification.

* `sans` - (Optional, String, ForceNew) Specifies additional domain names bound to multi-domain type certificates.
  This value only needs to be set when the purchased certificate is **multiple domain names** and there is a quota for
  adding additional domain names. Multiple domain names need to be separated by ";".
  For example, `www.example.com;www.example1.com;www.example2.com`.

  Changing this parameter will create a new resource.

* `csr` - (Optional, String, ForceNew) Specifies the certificate CSR string, which must match the domain name.

  Changing this parameter will create a new resource.

* `company_name` - (Optional, String, ForceNew) Specifies the company name. The valid length is limited from `0` to `63`.
  This field is required for OV and EV type certificates.

  Changing this parameter will create a new resource.

* `company_unit` - (Optional, String, ForceNew) Specifies the department name. The valid length is limited from `0` to `63`.

  Changing this parameter will create a new resource.

* `company_province` - (Optional, String, ForceNew) Specifies the province where the company is located. The valid length
  is limited from `0` to `63`. This field is required for OV and EV type certificates.

  Changing this parameter will create a new resource.

* `company_city` - (Optional, String, ForceNew) Specifies the city where the company is located. The valid length
  is limited from `0` to `63`. This field is required for OV and EV type certificates.

  Changing this parameter will create a new resource.

* `country` - (Optional, String, ForceNew) Specifies the country code. Must comply with the regular pattern `[A-Za-z]{2}`.
  This field is required for OV and EV type certificates.

  Changing this parameter will create a new resource.

* `contact_name` - (Optional, String, ForceNew) Specifies the technical contact name. The valid length is limited from
  `0` to `63`.

  Changing this parameter will create a new resource.

* `contact_phone` - (Optional, String, ForceNew) Specifies the technical contact phone number.

  Changing this parameter will create a new resource.

* `contact_email` - (Optional, String, ForceNew) Specifies the technical contact email.

  Changing this parameter will create a new resource.

* `auto_dns_auth` - (Optional, Bool, ForceNew) Specifies whether to push DNS verification information to HuaweiCloud
  resolution service.

  Changing this parameter will create a new resource.

* `key_algorithm` - (Optional, String, ForceNew) Specifies the key algorithm. Defaults to **RSA_2048**.

  Changing this parameter will create a new resource.

* `ca_hash_algorithm` - (Optional, String, ForceNew) Specifies the signature algorithm. This field is required for Geo
  OV certificate. Valid values are **DEFAULT** and **SHA-256**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
