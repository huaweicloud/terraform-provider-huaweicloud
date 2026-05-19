---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_ssl_cert_download"
description: |-
  Use this data source to obtain the SSL certificate download link for a specific DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_ssl_cert_download

Use this data source to obtain the SSL certificate download link for a specific DCS instance within HuaweiCloud.

> **NOTE:** This interface is currently only available for Redis 6.0/7.0 basic edition instances.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_ssl_cert_download" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to obtain the SSL certificate. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `file_name` - The file name of the SSL certificate.

* `link` - The download link of the SSL certificate.

* `bucket_name` - The OBS bucket name where the SSL certificate is stored.
