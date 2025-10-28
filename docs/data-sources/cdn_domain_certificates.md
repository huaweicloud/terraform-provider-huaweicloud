---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_certificates"
description: |-
  Use this data source to get the list of domains bound to HTTPS certificate of CDN within HuaweiCloud.
---

# huaweicloud_cdn_domain_certificates

Use this data source to get the list of domains bound to HTTPS certificate of CDN within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_cdn_domain_certificates" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the domain.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all** indicates all projects. This parameter is mandatory when you use an IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domain_certificates` - The list of certificates that are associated with the queried domain.
  The [domain_certificates](#cdn_domain_certificates) structure is documented below.

<a name="cdn_domain_certificates"></a>
The `domain_certificates` block supports:

* `domain_id` - The ID of the domain.

* `domain_name` - The name of the domain.

* `certificate_name` - The name of the certificate.

* `certificate_body` - The content of the certificate.

* `certificate_source` - The source type of the certificate.
  + **1**: Huawei-managed certificate.
  + **0**: Your own certificate.

* `expire_at` - The expiration time, in RFC3339 format.

* `https_status` - Whether the HTTPS certificate is enabled.
  + **0**: Do not enable HTTPS certificates.
  + **1**: Enable HTTPS acceleration and protocol follow back to origin.
  + **2**: Enable HTTPS acceleration and HTTP back to origin.

* `force_redirect_https` - Whether client requests are forced to be redirected.
  + **0**: Client requests will not be forced to redirect.
  + **1**: Client requests will be forced to redirect.
  + **2**: Client requests will be forced to jump to HTTP.

* `http2_enabled` - Whether the HTTP/2 protocol is used.
  + **0**: Not use HTTP2.0.
  + **1**: Use HTTP2.0.
