---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_segment_support_masks"
description: |-
  Use this dataSource to get the list of supported masks for global EIP segment.
---

# huaweicloud_global_eip_segment_support_masks

Use this dataSource to get the list of supported masks for global EIP segment.

## Example Usage

```hcl
data "huaweicloud_global_eip_segment_support_masks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `page_reverse` - (Optional, String) Specifies the page direction.
  The valid values are as follows:
  + **true**: means the previous page.  
  + **false**: means the next page.

* `fields` - (Optional, List) Specifies the fields to return.
  Supported values include **id**, **ip_version**, **mask**, **created_at**, and **updated_at**.

* `sort_key` - (Optional, String) Specifies the sort fields.

* `sort_dir` - (Optional, String) Specifies the sort directions.
  Valid values are **asc** and **desc**.

* `mask_ids` - (Optional, List) Specifies the mask IDs to filter.

* `ip_version` - (Optional, List) Specifies the IP versions to filter.

* `mask` - (Optional, Int) Specifies the mask length to filter. The value ranges from `1` to `128`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_masks` - The list of supported global EIP segment masks.

  The [support_masks](#support_masks_struct) structure is documented below.

<a name="support_masks_struct"></a>
The `support_masks` block supports:

* `id` - The mask record ID.

* `ip_version` - The IP version. The range of values is `4` and `6`

* `mask` - The mask length.

* `created_at` - The creation time.

* `updated_at` - The update time.
