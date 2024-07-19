---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificates"
description: |-
  Use this data source to get the list of CCM SSL certificates.
---

# huaweicloud_ccm_certificates

Use this data source to get the list of CCM SSL certificates.

## Example Usage

```hcl
data "huaweicloud_ccm_certificates" "test" {
  status = "ALL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the certificate name.

* `status` - (Optional, String) Specifies the certificate status.  
  The options are as follows:
  + **ALL**: All certificate status.
  + **PAID**: The certificate has been paid and needs to be applied for from the CA.
  + **ISSUED**: The certificate has been issued.
  + **CHECKING**: The certificate application is being reviewed.
  + **CANCELCHECKING**: The certificate application cancellation is being reviewed.
  + **UNPASSED**: The certificate application fails.
  + **EXPIRED**: The certificate has expired.
  + **REVOKING**: The certificate revocation application is being reviewed.
  + **REVOKED**: The certificate has been revoked.
  + **UPLOAD**: The certificate is being managed.
  + **CHECKING_ORG**: The organization verification is to be completed.
  + **ISSUING**: The certificate is to be issued.
  + **SUPPLEMENTCHECKING**: Additional domain names to be added for a multi-domain certificate are being reviewed.
  
  Defaults to **ALL**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the CCM certificates belong.
  This field is only valid for enterprise users. For enterprise users, all resources with permission will be queried
  when this field is not specified.

* `deploy_support` - (Optional, Bool) Specifies whether to query only certificates that support deployment.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - The certificate details list.
  The [certificates](#certificates_Certificate) structure is documented below.

<a name="certificates_Certificate"></a>
The `certificates` block supports:

* `id` - The certificate ID.

* `name` - The certificate name.

* `domain` - The domain name associated with the certificate.

* `sans` - The additional domain name associated with the certificate.

* `signature_algorithm` - The signature algorithm.

* `deploy_support` - Whether to support deployment.

* `type` - The certificate type. The value can be: **DV_SSL_CERT**, **DV_SSL_CERT_BASIC**, **EV_SSL_CERT**,
  **EV_SSL_CERT_PRO**, **OV_SSL_CERT**, **OV_SSL_CERT_PRO**.

* `brand` - The certificate authority. The value can be: **GLOBALSIGN**, **SYMANTEC**, **GEOTRUST**, **CFCA**.

* `expire_time` - The certificate expiration time.

* `domain_type` - The domain name type. The options are as follows:
  + **SINGLE_DOMAIN**: Single domain names.
  + **WILDCARD**: Wildcard domain names.
  + **MULTI_DOMAIN**: Multiple domain names.

* `validity_period` - The certificate validity period, in months.

* `status` - The certificate status.

* `domain_count` - The number of domain names that can be associated with the certificate.

* `wildcard_count` - The number of wildcard domain names that can be associated with the certificate.

* `description` - The certificate description.

* `enterprise_project_id` - The enterprise project ID.
