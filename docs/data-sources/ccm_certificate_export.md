---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_export"
description: |-
  Use this data source to export an SSL certificate within HuaweiCloud.
---

# huaweicloud_ccm_certificate_export

Use this data source to export an SSL certificate within HuaweiCloud.

## Example Usage

```hcl
variable "certificate_id" {}

data "huaweicloud_ccm_certificate_export" "test" {
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `certificate_id` - (Required, String) Specifies the SSL certificate ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enc_private_key` - The encryption certificate private key. This attribute is only meaningful in the state secret certificate.

* `entire_certificate` - The certificate content and certificate chain.

* `certificate` - The certificate content.

* `certificate_chain` - The certificate chain.

* `private_key` - The private key of the certificate.

* `enc_certificate` - The encryption certificate content. This attribute is only meaningful in the state secret certificate.
