---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_certificates"
description: |-
  Use this data source to get a list of WAF certificates within HuaweiCloud.
---

# huaweicloud_waf_certificates

Use this data source to get a list of WAF certificates within HuaweiCloud.

## Example Usage

```hcl
variable enterprise_project_id {}

data "huaweicloud_waf_certificates" "test" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of WAF certificate.
  For enterprise users, if omitted, default enterprise project will be used.

* `name` - (Optional, String) Specifies the name of certificate. The value is case-sensitive and supports fuzzy matching.

* `host` - (Optional, Bool) Specifies whether to obtain the domain name for which the certificate is used.
  + **true**: Obtain the certificates that have been used for domain names.
  + **false**: Obtain the certificates that have not been used for any domain names.

  Defaults to **false**.

* `expiration_status` - (Optional, String) Specifies the certificate expiration status. The options are as follows:
  + `0`: Not expired;
  + `1`: Expired;
  + `2`: Expired soon (The certificate will expire in one month.)

  -> If this field is not configured, all certificates that meet the expired status will be found.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - The certificate list.
  The [certificates](#certificates_struct) structure is documented below.

<a name="certificates_struct"></a>
The `certificates` block supports:

* `id` - The certificate ID.

* `name` - The certificate name.

* `expiration_status` - The certificate expiration status.

* `created_at` - The time when the certificate was uploaded, in RFC3339 format.

* `expired_at` - The time when the certificate expires, in RFC3339 format.

* `bind_host` - The domain information bound to the certificate.
  The [bind_host](#items_bind_host_struct) structure is documented below.

<a name="items_bind_host_struct"></a>
The `bind_host` block supports:

* `id` - The domain ID.

* `domain` - The domain name.

* `mode` - The special domain pattern.

* `waf_type` - The deployment mode of WAF instance that is used for the domain name. The value can be **cloud** for
  cloud WAF or **premium** for dedicated WAF instances.
