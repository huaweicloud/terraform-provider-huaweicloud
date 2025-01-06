---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_ssl_cert_download_link"
description: |-
  Use this data source to get the address for downloading the SSL certificate of a GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_ssl_cert_download_link

Use this data source to get the address for downloading the SSL certificate of a GaussDB OpenGauss instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_ssl_cert_download_link" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `download_link` - Indicates the download address of the SSL certificate.
