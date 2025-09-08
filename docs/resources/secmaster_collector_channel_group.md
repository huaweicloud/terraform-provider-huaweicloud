---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_channel_group"
description: |-
  Manages a collector channel group resource within HuaweiCloud.
---

# huaweicloud_secmaster_collector_channel_group

Manages a collector channel group resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "group_name" {}

resource "huaweicloud_secmaster_collector_channel_group" "test" {
  workspace_id = var.workspace_id
  name         = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the collector channel
  group belongs.

* `name` - (Required, String) Specifies the name of the collector channel group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The collector channel group can be imported using the `workspace_id` and their `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_collector_channel_group.test <workspace_id>/<name>
```
