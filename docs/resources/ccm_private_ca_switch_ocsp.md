---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_switch_ocsp"
description: |-
  Manages a CCM private CA switch OCSP resource within HuaweiCloud.
---

# huaweicloud_ccm_private_ca_switch_ocsp

Manages a CCM private CA switch OCSP resource within HuaweiCloud.

-> Deleting this resource will not recover the private CA OCSP switch status, but will only remove the
  resource information from the tfstate file.

## Example Usage

```hcl
variable "ca_id" {}

resource "huaweicloud_ccm_private_ca_switch_ocsp" "test" {
  ca_id       = var.ca_id
  ocsp_switch = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `ca_id` - (Required, String, NonUpdatable) Specifies the ID of the private CA you want to switch OCSP.

* `ocsp_switch` - (Required, Bool) Specifies whether to enable or disable OCSP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the private CA ID).
