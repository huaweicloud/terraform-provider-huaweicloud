---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_recycle_bin_policy"
description: |-
  Manages an EVS update recycle bin policy resource within HuaweiCloud.
---

# huaweicloud_evs_recycle_bin_policy

Manages an EVS update recycle bin policy resource within HuaweiCloud.

-> This resource is a one-time action resource using to update EVS recycle bin policy. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_evs_recycle_bin_policy" "test" {
  switch         = true
  threshold_time = 7
  keep_time      = 7
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `switch` - (Optional, Bool) Specifies whether the recycle bin switch.  
  The valid values are as follows:
  + **true**: Open recycle bin.
  + **false**: Close recycle bin.

* `threshold_time` - (Optional, Int) Specifies the threshold time for the recycle bin, and how many days after the cloud
  disk is created, it will be deleted before being placed in the recycle bin.  
  The valid value range is `1` to `1,000`. The default value is `7`.

* `keep_time` - (Optional, Int) Specifies the storage period (days) of the cloud disk in the designated recycle bin, and
  how many days it will take for the cloud disk to be completely deleted after entering the recycle bin.  
  The valid value range is `1` to `365`. The default value is `7`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
