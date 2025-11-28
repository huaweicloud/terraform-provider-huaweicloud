---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_application_batch_action"
description: |-
  Use this resource to batch enable or disable applications of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_application_batch_action

Use this resource to batch enable or disable applications of the Workspace APP within HuaweiCloud.

-> This resource is a one-time action resource used to batch enable or disable applications. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

```hcl
variable "app_group_id" {}
variable "application_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_application_batch_action" "test" {
  app_group_id    = var.app_group_id
  action          = "disable"
  application_ids = var.application_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the applications to be operated are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_group_id` - (Required, String, NonUpdatable) Specifies the ID of the application group.

* `action` - (Required, String, NonUpdatable) Specifies the type of the action.  
  The valid values are as follows:
  + **enable**
  + **disable**

* `application_ids` - (Required, List, NonUpdatable) Specifies the list of application IDs to be operated.  

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
