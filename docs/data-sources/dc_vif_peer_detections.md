---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_vif_peer_detections"
description:
  Use this data source to get the list of DC virtual interface peer detections.
---

# huaweicloud_dc_vif_peer_detections

Use this data source to get the list of DC virtual interface peer detections.

## Example Usage

```hcl
variable "vif_peer_id" {}

data "huaweicloud_dc_vif_peer_detections" "test" {
  vif_peer_id = var.vif_peer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vif_peer_id` - (Required, String) Specifies the ID of the virtual interface peer.

* `sort_key` - (Optional, String) Specifies the sort key.

* `sort_dir` - (Optional, List) Specifies the list of sort dir. Value options: **desc**, **asc**. Defaults to **asc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vif_peer_detections` - Indicates the list of virtual interface peer detection.
  The [vif_peer_detections](#vif_peer_detections_struct) structure is documented below.

<a name="vif_peer_detections_struct"></a>
The `vif_peer_detections` block supports:

* `id` - The virtual interface peer detection ID.

* `vif_peer_id` - Indicates the virtual interface peer ID.

* `project_id` - Indicates the project ID.

* `domain_id` - Indicates the domain ID.

* `status` - Indicates the status of the virtual interface peer detection.

* `start_time` - Indicates the start time of the virtual interface peer detection.

* `end_time` - Indicates the end time of the virtual interface peer detection.

* `loss_ratio` - Indicates the loss ratio.
