---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_metadata"
description: |-
  Use this resource to manage EVS volume metadata within HuaweiCloud.
---

# huaweicloud_evs_volume_metadata

Use this resource to manage EVS volume metadata within HuaweiCloud.

## Example Usage

```hcl
variable "volume_id" {}

resource "huaweicloud_evs_volume_metadata" "test" {
  volume_id = var.volume_id

  metadata {
    key1 = "value1"
    key2 = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `volume_id` - (Required, String, NonUpdatable) Specifies the disk ID.

* `metadata` - (Required, Map) Specifies the user-defined metadata key-value pair.

  -> When updating the `metadata` parameter, all existing key-value pairs will be overwritten.
    <br/>`key` or `value` under metadata can contain no more than `255` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `volume_id`.
