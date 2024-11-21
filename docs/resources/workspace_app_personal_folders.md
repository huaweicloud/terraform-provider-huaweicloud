---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_personal_folders"
description: |-
  Using this resource to assign or manages the personal folders under NAS storage of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_personal_folders

Using this resource to assign or manages the personal folders under NAS storage of Workspace APP within HuaweiCloud.

## Example Usage

### Create a personal folder for each user with the same permission policy

```hcl
variable "workspace_app_nas_storage_id" {}
variable "workspace_app_storage_policy_id" {}
variable "workspace_user_names" {
  type = list(string)
}

resource "huaweicloud_workspace_app_personal_folders" "test" {
  storage_id = var.workspace_app_nas_storage_id

  dynamic "assignments" {
    for_each = var.workspace_user_names

    content {
      policy_statement_id = var.workspace_app_storage_policy_id
      attach              = assignments.value
      attach_type         = "USER"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the personal folders are located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `storage_id` - (Required, String, ForceNew) Specifies the NAS storage ID to which the personal folders belong.  
  Change this parameter will create a new resource.

* `assignments` - (Required, List, ForceNew) Specifies the assignment configuration of personal folders.  
  The [assignments](#workspace_app_personal_folder_assignment) structure is documented below.

<a name="workspace_app_personal_folder_assignment"></a>
The `assignments` block supports:

* `policy_statement_id` - (Required, String, ForceNew) Specifies the ID of the storage permission policy.  
  Change this parameter will create a new resource.

* `attach` - (Required, String, ForceNew) Specifies the object name of personal folder assignment.  
  Change this parameter will create a new resource.

* `attach_type` - (Optional, String, ForceNew) Specifies the type of personal folder assignment.  
  The valid value is **USER** (also default value).  
  Change this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `storage_id`.

## Import

Personal folders can be imported using resource `id`, also ID of NAS storage to which personal folders belong, e.g.

```bash
$ terraform import huaweicloud_workspace_app_personal_folders.test <id>
```

If the NAS storage ID is unknow, the NAS storage name can be used as an alternative to ID.  
The NAS storage name can be obtained through the console or data source (`huaweicloud_workspace_app_nas_storages`).

```bash
$ terraform import huaweicloud_workspace_app_personal_folders.test <storage_name>
```
