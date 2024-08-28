---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_ca_restore"
description: |-
  Manages a CCM private CA restore resource within HuaweiCloud.
---

# huaweicloud_ccm_private_ca_restore

Manages a CCM private CA restore resource within HuaweiCloud.

-> 1. This resource is only applicable to private CAs in the **DELETED** state.
<br/>2. The current resource is a one-time resource, and destroying this resource will not affect the result.

## Example Usage

```hcl
variable "ca_id" {}

resource "huaweicloud_ccm_private_ca_restore" "test" {
  ca_id = var.ca_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `ca_id` - (Required, String, ForceNew) Specifies the ID of the private CA you want to restore. The specified private
  CA status must be **DELETED**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
