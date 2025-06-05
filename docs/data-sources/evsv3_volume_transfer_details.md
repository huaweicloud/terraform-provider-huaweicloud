---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evsv3_volume_transfer_details"
description: |-
  Use this data source to get the list of EVS volume transfer details (V3) within HuaweiCloud.
---

# huaweicloud_evsv3_volume_transfer_details

Use this data source to get the list of EVS volume transfer details (V3) within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evsv3_volume_transfer_details" "test" {}
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

* `created_at` - The time when the transfer was created.

* `links` - The links of the cloud disk transfer record.
  The [links](#links_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `href` - The corresponding shortcut link.

* `rel` - The shortcut link marker name.
