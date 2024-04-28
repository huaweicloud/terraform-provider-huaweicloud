---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_maintainwindow"
description: ""
---

# huaweicloud_dcs_maintainwindow

Use this data source to get the ID of an available DCS maintenance windows.

## Example Usage

```hcl
data "huaweicloud_dcs_maintainwindow" "maintainwindow1" {
  seq = 1
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dcs maintenance windows. If omitted, the provider-level
  region will be used.

* `seq` - (Optional, Int) Specifies the sequential number of a maintenance time window.

* `begin` - (Optional, String) Specifies the time at which a maintenance time window starts.

* `end` - (Optional, String) Specifies the time at which a maintenance time window ends.

* `default` - (Optional, Bool) Specifies whether a maintenance time window is set to the default time segment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.
