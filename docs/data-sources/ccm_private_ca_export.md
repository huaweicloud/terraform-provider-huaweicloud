---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_export"
description: |-
  Use this data source to export a private CA within HuaweiCloud.
---

# huaweicloud_ccm_private_ca_export

Use this data source to export a private CA within HuaweiCloud.

-> Only CAs in `ACTIVED`, `DISABLED` or `EXPIRED` status support exporting operation.

## Example Usage

```hcl
variable "ca_id" {}

data "huaweicloud_ccm_private_ca_export" "test" {
  ca_id = var.ca_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ca_id` - (Required, String) Specifies the ID of the CA you want to export.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificate` - The certificate content.

* `certificate_chain` - The content of the certificate chain.
