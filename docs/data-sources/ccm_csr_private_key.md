---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_csr_private_key"
description: |-
  Use this data source to get the CCM SSL CSR private key.
---

# huaweicloud_ccm_csr_private_key

Use this data source to get the CCM SSL CSR private key.

## Example Usage

```hcl
variable "csr_id" {}

data "huaweicloud_ccm_csr_private_key" "test" {
  csr_id = var.csr_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `csr_id` - (Required, String) Specifies the CSR ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `private_key` - The private key.
