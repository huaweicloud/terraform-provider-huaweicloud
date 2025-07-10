---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_unsubscribe_prepaid_volume"
description: |-
  Manages a resource to unsubscribe the prepaid EVS volume within HuaweiCloud.
---

# huaweicloud_evs_unsubscribe_prepaid_volume

Manages a resource to unsubscribe the prepaid EVS volume within HuaweiCloud.

-> 1. It cannot be used to unsubscribe from system disks and bootable disks. They must be unsubscribed from together
with their servers.<br>2. A maximum of 60 disks can be unsubscribed from at the same time using this API.<br>3. The
current resource is a one-time resource, and destroying this resource will not recover the unsubscribe volume,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "volume_ids" {}

resource "huaweicloud_evs_unsubscribe_prepaid_volume" "example" {
  volume_ids = var.volume_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `volume_ids` - (Required, List, NonUpdatable) Specifies the volume IDs. A maximum of `60` volume IDs can be
  configured.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
