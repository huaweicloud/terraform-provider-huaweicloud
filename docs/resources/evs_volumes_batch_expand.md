---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volumes_batch_expand"
description: |-
  Manages an EVS volumes expand resource within HuaweiCloud.
---

# huaweicloud_evs_volumes_batch_expand

Manages an EVS volumes expand resource within HuaweiCloud.

-> The current resource is a one-time action resource using to expand volumes. Deleting this resource will not reset
the expanded volumes, but will only remove the resource information from the tfstate file.  
Using this resource may cause incompatible changes to other resources that contain volume size fields. Please use
`lifecycle.ignore_changes` properly to handle unexpected changes.  
If the status of the to-be-expanded disk is available, there are no restrictions. If the status of the to-be-expanded
disk is in-use, the restrictions are as follows:
<br/>1. A shared disk cannot be expanded, which means that the value of multiattach must be false.
<br/>2. The status of the server to which the disk attached must be **ACTIVE**, **PAUSED**, **SUSPENDED**, or **SHUTOFF**.

## Example Usage

```hcl
variable "id" {}
variable "new_size" {}
variable "is_auto_pay" {
  type = bool
}

resource "huaweicloud_evs_volumes_batch_expand" "test" {
  volumes    {
    id       = var.id
    new_size = var.new_size
  }

  is_auto_pay = var.is_auto_pay
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `volumes` - (Required, List, NonUpdatable) Specifies the to-be-expanded volume list.
  The [volumes](#volumes_struct) structure is documented below.

* `is_auto_pay` - (Optional, Bool, NonUpdatable) Specifies whether to pay immediately. This parameter is valid only
  when the disk in prepaid mode. Defaults to **false**. Possible values are:
  + **true**: An order is immediately paid from the account balance.
  + **false**: An order is not paid immediately after being created.

<a name="volumes_struct"></a>
The `volumes` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the volume ID.

* `new_size` - (Required, Int, NonUpdatable) Specifies the new size of the to-be-expanded volume, in GiB.
  Must be greater than the current size. The maximum disk size: Data disk: `32,768` GiB, System disk: `1,024` GiB

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
