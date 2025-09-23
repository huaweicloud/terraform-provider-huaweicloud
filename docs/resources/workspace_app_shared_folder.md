---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_shared_folder"
description: |-
  Manages a shared folder resource of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_shared_folder

Manages a shared folder resource of Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "nas_storage_id" {}
variable "shared_folder_name" {}

resource "huaweicloud_workspace_app_shared_folder" "test" {
  storage_id = var.nas_storage_id
  name       = var.shared_folder_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the shared folder is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `storage_id` - (Required, String, ForceNew) Specifies the NAS storage ID to which the shared folder belongs.  
  Change this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the shared folder.  
  The valid length is limited from `1` to `32`, only letters, digits, spaces, underscores (_) and hyphens (-) are
  allowed, but cannot be all spaces or start with the space.  
  Change this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The shared folder ID.

* `delimiter` - The delimiter that the shared folder path used.

* `path` - The path of the shared folder, the usual format is: **{root_path}{delimiter}{folder_name}{delimiter}**,
  such as **shares/xxx/**.  

## Import

Shared folders can be imported using their related `storage_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_workspace_app_shared_folder.test <storage_id>/<id>
```

If the NAS storage ID or shared folder ID is unknow, its corresponding name can be used as an alternative to ID.  
The NAS storage name can be obtained through the console or data source (`huaweicloud_workspace_app_nas_storages`).

```bash
$ terraform import huaweicloud_workspace_app_shared_folder.test <storage_name>/<name>
```
