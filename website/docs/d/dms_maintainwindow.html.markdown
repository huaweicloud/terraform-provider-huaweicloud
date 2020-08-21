---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_maintainwindow"
sidebar_current: "docs-huaweicloud-datasource-dms-maintainwindow"
description: |-
  Get information on an HuaweiCloud dms maintainwindow.
---

# huaweicloud\_dms\_maintainwindow

Use this data source to get the ID of an available HuaweiCloud dms maintainwindow.
This is an alternative to `huaweicloud_dms_maintainwindow_v1`

## Example Usage

```hcl

data "huaweicloud_dms_maintainwindow" "maintainwindow1" {
  seq = 1
}

```

## Argument Reference

* `seq` - (Required) Indicates the sequential number of a maintenance time window.

* `begin` - (Optional) Indicates the time at which a maintenance time window starts.

* `end` - (Required) Indicates the time at which a maintenance time window ends.

* `default` - (Required) Indicates whether a maintenance time window is set to the default time segment.

## Attributes Reference

`id` is set to the ID of the found maintainwindow. In addition, the following attributes
are exported:

* `begin` - See Argument Reference above.
* `end` - See Argument Reference above.
* `default` - See Argument Reference above.
