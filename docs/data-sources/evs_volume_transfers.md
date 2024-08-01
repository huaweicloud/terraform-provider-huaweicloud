---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_transfers"
description: |-
  Use this data source to get the list of EVS volume transfers within HuaweiCloud.
---

# huaweicloud_evs_volume_transfers

Use this data source to get the list of EVS volume transfers within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_volume_transfers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `transfers` - The list of volume transfers.

  The [transfers](#transfers_struct) structure is documented below.

<a name="transfers_struct"></a>
The `transfers` block supports:

* `id` - The volume transfer ID.

* `name` - The volume transfer name.

* `volume_id` - The volume ID.
