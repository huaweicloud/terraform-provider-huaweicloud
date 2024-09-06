---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_event_report"
description: |-
  Manages a CES event report resource within HuaweiCloud.
---

# huaweicloud_ces_event_report

Manages a CES event report resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "name" {}
variable "time" {}
variable "source" {}
variable "content" {}
variable "resource_id" {}
variable "resource_name" {}
variable "user" {}

resource "huaweicloud_ces_event_report" "test" {
  name   = var.name
  source = var.source
  time   = var.time

  detail {
    state         = "normal"
    level         = "Major"
    content       = var.content
    resource_id   = var.resource_id
    resource_name = var.resource_name
    user          = var.user
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the event name.

* `source` - (Required, String, NonUpdatable) Specifies the event source.

* `time` - (Required, String, NonUpdatable) Specifies the occurrence time of the event.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `detail` - (Required, List, NonUpdatable) Specifies the detail of the CES event.
  The [detail](#Detail) structure is documented below.

<a name="Detail"></a>
The `detail` block supports:

* `state` - (Required, String, NonUpdatable) Specifies the event status.
  The value can be **normal**, **warning**, or **incident**.

* `level` - (Required, String, NonUpdatable) Specifies the event level.
  The value can be **Critical**, **Major**, **Minor**, or **Info**.

* `type` - (Optional, String, NonUpdatable) Specifies the event type.
  The value can only be **EVENT.CUSTOM**.

* `content` - (Optional, String, NonUpdatable) Specifies the event content.

* `group_id` - (Optional, String, NonUpdatable) Specifies the group that the event belongs to.

* `resource_id` - (Optional, String, NonUpdatable) Specifies the resource ID.

* `resource_name` - (Optional, String, NonUpdatable) Specifies the resource name.

* `user` - (Optional, String, NonUpdatable) Specifies the event user.

* `dimensions` - (Optional, List, NonUpdatable) Specifies the resource dimensions.
  The [dimensions](#DetailDimensions) structure is documented below.

<a name="DetailDimensions"></a>
The `dimensions` block supports:

* `name` - (Required, String, NonUpdatable) The resource dimension name.

* `value` - (Required, String, NonUpdatable) The resource dimension value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
