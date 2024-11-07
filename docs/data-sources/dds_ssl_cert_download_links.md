---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_ssl_cert_download_links"
description: |-
  Use this data source to get the list of DDS instance ssl cert download links.
---

# huaweicloud_dds_ssl_cert_download_links

Use this data source to get the list of DDS instance ssl cert download links.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_ssl_cert_download_links" "test"{
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

* `certs` - Indicates the certificate list.

  The [certs](#certs_struct) structure is documented below.

<a name="certs_struct"></a>
The `certs` block supports:

* `download_link` - Indicates the certificate download link.

* `category` - Indicates the certificate type.
