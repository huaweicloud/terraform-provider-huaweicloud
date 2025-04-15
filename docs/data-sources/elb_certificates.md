---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_certificates"
description: |-
  Use this data source to get the list of ELB certificates.
---

# huaweicloud_elb_certificates

Use this data source to get the list of ELB certificates.

## Example Usage

```hcl
data "huaweicloud_elb_certificates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `certificate_id` - (Optional, List) Specifies a certificate ID.
  Multiple IDs can be queried.

* `name` - (Optional, List) Specifies the certificate name.
  Multiple names can be queried.

* `description` - (Optional, List) Specifies the certificate description.
  Multiple descriptions can be queried.

* `domain` - (Optional, List) Specifies the domain names used by the server certificate.
  Multiple domain names can be queried.

* `type` - (Optional, List) Specifies the certificate type.
  Value options:
  + **server** indicates server certificates.
  + **client** indicates CA certificates.
  + **server_sm**: indicates the server SM certificate.
  Multiple types can be queried.

* `common_name` - (Optional, List) Specifies the primary domain name of the certificate.
  Multiple values can be queried.

* `fingerprint` - (Optional, List) Specifies the fingerprint of the certificate.
  Multiple values can be queried.

* `scm_certificate_id` - (Optional, List) Specifies the SSL certificate ID.
  Multiple IDs can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - Indicates the list of certificates.

  The [certificates](#certificates_struct) structure is documented below.

<a name="certificates_struct"></a>
The `certificates` block supports:

* `id` - Indicates the certificate ID.

* `name` - Indicates  the certificate name.

* `type` - Indicates the certificate type.

* `description` - Indicates  the certificate description.

* `subject_alternative_names` - Indicates all the domain names of the certificate.

* `certificate` - Indicates the certificate content.

* `fingerprint` - Indicates the fingerprint of the certificate.

* `domain` - Indicates the domain names used by the server certificate.

* `private_key` - Indicates the private key of the certificate used by HTTPS listeners.

* `common_name` - Indicates the primary domain name of the certificate.

* `enc_certificate` - Indicates the body of the SM encryption certificate required by HTTPS listeners.

* `scm_certificate_id` - Indicates the SSL certificate ID.

* `enc_private_key` - Indicates the private key of the SM encryption certificate required by HTTPS listeners.

* `expire_time` - Indicates the time when the certificate expires.

* `created_at` - Indicates the time when the certificate was created.

* `updated_at` - Indicates the time when the certificate was updated.
