---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_adb_command"
description: |-
  Manages a CPH adb command resource within HuaweiCloud.
---

# huaweicloud_cph_adb_command

Manages a CPH adb command resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "command" {}
variable "content" {}
variable "phone_ids" {}

resource "huaweicloud_cph_adb_command" "test" {
  command   = var.command
  content   = var.content
  phone_ids = var.phone_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `command` - (Required, String, NonUpdatable) Specifies the ADB command. The value can be **push**, **install**,
  **install-multiple**, **uninstall**, or **shell**.

* `content` - (Required, String, NonUpdatable) Specifies the content.
  + When the `command` is **push**, **install**, **install-multiple** or **uninstall**, it indicates OBS object path.
    For example, **obs://obs-bucket-name/obs-file-path/file.tar**.
  + When the `command` is **shell**, it indicates the command to be executed. For example, **ls -l /system**.

* `phone_ids` - (Optional, List, NonUpdatable) Specifies the IDs of the CPH phone.

* `server_ids` - (Optional, List, NonUpdatable) Specifies the IDs of CPH server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
