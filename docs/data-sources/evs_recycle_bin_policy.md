---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_recycle_bin_policy"
description: |-
  Use this data source to get the EVS recycle bin policy within HuaweiCloud.
---

# huaweicloud_evs_recycle_bin_policy

Use this data source to get the EVS recycle bin policy within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_recycle_bin_policy" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `switch` - Whether the recycle bin switch.  
  The valid values are as follows:
  + **true**: Open recycle bin.
  + **false**: Close recycle bin.

* `threshold_time` - The threshold time for the recycle bin, and how many days after the cloud
  disk is created, it will be deleted before being placed in the recycle bin.  
  The valid value range is `1` to `1,000`. The default value is `7`.

* `keep_time` - The storage period (days) of the cloud disk in the designated recycle bin, and
  how many days it will take for the cloud disk to be completely deleted after entering the recycle bin.  
  The valid value range is `1` to `365`. The default value is `7`.
