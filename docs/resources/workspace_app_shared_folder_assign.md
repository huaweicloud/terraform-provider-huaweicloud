---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_shared_folder_assign"
description: |-
  Use this resource to assign user access to the shared folder within HuaweiCloud.
---

# huaweicloud_workspace_app_shared_folder_assign

Use this resource to assign user access to the shared folder within HuaweiCloud.

-> This resource is only a one-time action resource for assignment of user access to a shared folder.
   Deleting this resource will not clear the corresponding request record, but will only remove
   the resource information from the tfstate file.

## Example Usage

```hcl
variable "storage_id" {}
variable "storage_claim_id" {}
variable "add_items" {
  type = list(object({
    policy_statement_id = string
    attach              = string
    attach_type         = string
  }))
}
variable "del_items" {
  type = list(object({
    attach      = string
    attach_type = string
  }))
}

resource "huaweicloud_workspace_app_shared_folder_assign" "test" {
  storage_id       = var.storage_id
  storage_claim_id = var.storage_claim_id

  dynamic "add_items" {
    for_each = var.add_items

    content {
      policy_statement_id = add_items.value.policy_statement_id
      attach              = add_items.value.attach
      attach_type         = add_items.value.attach_type
    }
  }

  dynamic "del_items" {
    for_each = var.del_items

    content {
      attach      = del_items.value.attach
      attach_type = del_items.value.attach_type
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the WKS storage is located.  
  If omitted, the provider-level region will be used.  
  Change this parameter will create a new resource.

* `storage_id` - (Required, String, NonUpdatable) Specifies the WKS storage ID to which the shared folder belongs.

* `storage_claim_id` - (Required, String, NonUpdatable) Specifies the WKS storage directory claim ID.

* `add_items` - (Optional, List, NonUpdatable) Specifies the list of members to be added.  
  The [add_items](#workspace_shared_folders_access_add_items) structure is documented below.

* `del_items` - (Optional, List, NonUpdatable) Specifies the list of members to be removed.  
  The [del_items](#workspace_shared_folders_access_del_items) structure is documented below.

<a name="workspace_shared_folders_access_add_items"></a>
The `add_items` block supports:

* `policy_statement_id` - (Required, String) Specifies the policy ID.

* `attach` - (Required, String) Specifies the target.

* `attach_type` - (Required, String) Specifies the associated object type.  
  The valid values are as follows:
  + **USER** - User
  + **USER_GROUP** - User group

<a name="workspace_shared_folders_access_del_items"></a>
The `del_items` block supports:

* `attach` - (Required, String) Specifies the target.

* `attach_type` - (Required, String) Specifies the associated object type.  
  The valid values are as follows:
  + **USER** - User
  + **USER_GROUP** - User group

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
