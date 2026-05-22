---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_ssl_cert_download_links"
description: |-
  Use this data source to get the SSL certificate download link.
---

# huaweicloud_geminidb_ssl_cert_download_links

Use this data source to get the SSL certificate download link.

-> This data source does not support query the SSL certificate and private download addresses of CCM.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_ssl_cert_download_links" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certs` - The list of the certificate information.
  The [certs](#certs_struct) structure is documented below.

<a name="certs_struct"></a>
The `certs` block supports:

* `category` - The certificate type.
  + **international**: international certificate.

* `download_link` - The certificate download link.
