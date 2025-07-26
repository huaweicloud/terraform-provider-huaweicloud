---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_notification"
description: |-
  Use this resource to dispatch desktop message within HuaweiCloud.
---

# huaweicloud_workspace_desktop_notification

Use this resource to dispatch desktop message within HuaweiCloud.

-> This resource is only a one-time action resource for dispatch desktop message. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "desktop_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_desktop_notification" "test" {
  desktop_ids   = var.desktop_ids
  notifications = "terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktop is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `desktop_ids` - (Required, List, NonUpdatable) Specifies the list of desktop IDs.

* `notifications` - (Required, String, NonUpdatable) Specifies the message want to dispatch.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of notification dispatch task.  
  The valid values are as follows:
  + **SUCCESS**
  + **FAIL**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
