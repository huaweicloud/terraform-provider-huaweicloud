---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_transfer_accepter"
description: |-
  Manages an EVS volume transfer accepter resource within HuaweiCloud.
---

# huaweicloud_evs_volume_transfer_accepter

Manages an EVS volume transfer accepter resource within HuaweiCloud.

-> After successfully accepting the transfer of volume using this resource, the original volume transfer resource will
   no longer exist. Destroying resource does not change the current state of the accepter resource.

## Example Usage

```hcl
variable "transfer_id" {}
variable "auth_key" {}

resource "huaweicloud_evs_volume_transfer_accepter" "test" {
  transfer_id = var.transfer_id
  auth_key    = var.auth_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `transfer_id` - (Required, String, ForceNew) Specifies the ID of the volume transfer record.
  Changing this parameter will create a new resource.

* `auth_key` - (Required, String, ForceNew) Specifies the identity authentication key for volume transfer.
  When creating the volume transfer, the value of this field will be returned.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `transfer_id`.
