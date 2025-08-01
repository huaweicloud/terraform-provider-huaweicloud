---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_one_click_alarm_reset"
description: |-
  Manages a CES reset alarm rules for one service in one-click monitoring resource within HuaweiCloud.
---

# huaweicloud_ces_one_click_alarm_reset

Manages a CES reset alarm rules for one service in one-click monitoring resource within HuaweiCloud.

-> Deleting the reset alarm rules resource is not supported. The reset alarm rules for one service in one-click
monitoring resource is only removed from the state.

## Example Usage

```hcl
variable "one_click_alarm_id" {}

resource "huaweicloud_ces_one_click_alarm_reset" "test" {
  one_click_alarm_id = var.one_click_alarm_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `one_click_alarm_id` - (Required, String, NonUpdatable) Specifies the one-click monitoring ID for a service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
