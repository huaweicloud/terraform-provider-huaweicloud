---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_scm_certificates"
description: ""
---

# huaweicloud_scm_certificates

!> **WARNING:** It has been deprecated, use `huaweicloud_ccm_certificates` to get the CCM SSL certificates.

Use this data source to get the list of SCM certificates.

## Example Usage

```hcl
data "huaweicloud_scm_certificates" "test" {
  status = "ALL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Certificate name.

* `status` - (Optional, String) Certificate status.  
The options are as follows:
  - ALL: All certificate status.
  - PAID: The certificate has been paid and needs to be applied for from the CA.
  - ISSUED: The certificate has been issued.
  - CHECKING: The certificate application is being reviewed.
  - CANCELCHECKING: The certificate application cancellation is being reviewed.
  - UNPASSED: The certificate application fails.
  - EXPIRED: The certificate has expired.
  - REVOKING: The certificate revocation application is being reviewed.
  - REVOKED: The certificate has been revoked.
  - UPLOAD: The certificate is being managed.
  - CHECKING_ORG: The organization verification is to be completed.
  - ISSUING: The certificate is to be issued.
  - SUPPLEMENTCHECKING: Additional domain names to be added for a multi-domain certificate are being reviewed.
  
  Default: ALL

* `enterprise_project_id` - (Optional, String) The enterprise project id of the project.

* `deploy_support` - (Optional, Bool) Whether to support deployment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - Certificate list. For details, see Data structure of the Certificate field.
  The [Certificate](#certificates_Certificate) structure is documented below.

<a name="certificates_Certificate"></a>
The `Certificate` block supports:

* `id` - Certificate ID.

* `name` - Certificate name.

* `domain` - Domain name associated with the certificate.

* `sans` - Additional domain name associated with the certificate.

* `signature_algorithm` - Signature algorithm.

* `deploy_support` - Whether to support deployment.

* `type` - Certificate type.  
The value can be: **DV_SSL_CERT**, **DV_SSL_CERT_BASIC**, **EV_SSL_CERT**, **EV_SSL_CERT_PRO**, **OV_SSL_CERT**, **OV_SSL_CERT_PRO**.

* `brand` - Certificate authority.  
The value can be: **GLOBALSIGN**, **SYMANTEC**, **GEOTRUST**, **CFCA**.

* `expire_time` - Certificate expiration time.

* `domain_type` - Domain name type.  
The options are as follows:
  - SINGLE_DOMAIN: Single domain names
  - WILDCARD: Wildcard domain names
  - MULTI_DOMAIN: Multiple domain names

* `validity_period` - Certificate validity period, in months.

* `status` - Certificate status.  
The options are as follows:
  - ALL: All certificate status.
  - PAID: The certificate has been paid and needs to be applied for from the CA.
  - ISSUED: The certificate has been issued.
  - CHECKING: The certificate application is being reviewed.
  - CANCELCHECKING: The certificate application cancellation is being reviewed.
  - UNPASSED: The certificate application fails.
  - EXPIRED: The certificate has expired.
  - REVOKING: The certificate revocation application is being reviewed.
  - REVOKED: The certificate has been revoked.
  - UPLOAD: The certificate is being managed.
  - CHECKING_ORG: The organization verification is to be completed.
  - ISSUING: The certificate is to be issued.
  - SUPPLEMENTCHECKING: Additional domain names to be added for a multi-domain certificate are being reviewed.

* `domain_count` - Number of domain names that can be associated with the certificate.

* `wildcard_count` - Number of wildcard domain names that can be associated with the certificate.

* `description` - Certificate description.

* `enterprise_project_id` - The enterprise project id of the project.
