---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_vif_peer_detection"
description: |-
  Manages a DC vif peer detection resource within HuaweiCloud.
---

# huaweicloud_dc_vif_peer_detection

Manages a DC vif peer detection resource within HuaweiCloud.

## Example Usage

```hcl
variable "vif_peer_id" {}

resource "huaweicloud_dc_vif_peer_detection" "test" {
  vif_peer_id = var.vif_peer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `vif_peer_id` - (Required, String, NonUpdatable) Specifies the ID of the virtual interface peer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the virtual interface peer detection.

* `start_time` - Indicates the start time of the virtual interface peer detection.

* `end_time` - Indicates the end time of the virtual interface peer detection.

* `loss_ratio` - Indicates the loss ratio.

## Import

The DC vif peer detection resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dc_vif_peer_detection.test <id>
```
