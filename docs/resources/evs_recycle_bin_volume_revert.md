---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_recycle_bin_volume_revert"
description: |-
  Manages an EVS revert recycle bin volume resource within HuaweiCloud.
---

# huaweicloud_evs_recycle_bin_volume_revert

Manages an EVS revert recycle bin volume resource within HuaweiCloud.

-> This resource is a one-time action resource using to revert EVS recycle bin volume. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "volume_id" {}

resource "huaweicloud_evs_recycle_bin_volume_revert" "test" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the disk ID.
  For its values, can be obtained using `huaweicloud_evs_volumes` dataSource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `volume_id`.
