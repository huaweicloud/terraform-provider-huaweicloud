---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_visibility_batch_action"
description: |-
  Use this resource to batch set application visibility within HuaweiCloud.
---

# huaweicloud_workspace_application_visibility_batch_action

Use this resource to batch set application visibility within HuaweiCloud.

-> This resource is a one-time action resource used to batch set application visibility. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "application_visibility_action" {}
variable "operate_application_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_application_visibility_batch_action" "test" {
  action  = var.application_visibility_action
  app_ids = var.operate_application_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the applications to be operated are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the visibility action type.  
  Valid values are:
  + **enable**
  + **disable**

* `app_ids` - (Required, List, NonUpdatable) Specifies the list of application IDs to be operated.  
  Maximum of `50` applications are supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
