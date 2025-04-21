---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_ssl_cert_download_links"
description: |-
  Use this data source to get the list of RDS instance SSL certificates.
---

# huaweicloud_rds_ssl_cert_download_links

Use this data source to get the list of RDS instance SSL certificates.

## Example Usage

```hcl
data "huaweicloud_rds_ssl_cert_download_links" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cert_info_list` - Indicates the list of certificates.

  The [cert_info_list](#cert_info_list_struct) structure is documented below.

<a name="cert_info_list_struct"></a>
The `cert_info_list` block supports:

* `download_link` - Indicates the download link of certificate.

* `category` - Indicates the category of certificate.
  The value can be: **international**, **national**.
