---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_nas_storage"
description: |-
  Manages a NAS storage resource of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_nas_storage

Manages a NAS storage resource of Workspace APP within HuaweiCloud.

## Example Usage

### Create a NAS storage via SFS (3.0) file system

```hcl
variable "nas_storage_name" {}
variable "sfs_file_system_name" {}

resource "huaweicloud_workspace_app_nas_storage" "test" {
  name = var.nas_storage_name

  storage_metadata {
    storage_handle = var.sfs_file_system_name
    storage_class  = "sfs"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the NAS storage is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the NAS storage.  
  The valid length is limited from `1` to `128`, and allows visible characters or spaces, but cannot be all spaces.
  Change this parameter will create a new resource.

* `storage_metadata` - (Required, List, ForceNew) Specifies the metadata of the corresponding storage.  
  The [storage_metadata](#workspace_app_nas_storage_metadata) structure is documented below.

<a name="workspace_app_nas_storage_metadata"></a>
The `storage_metadata` block supports:

* `storage_handle` - (Required, String, ForceNew) Specifies the storage name.  
  Change this parameter will create a new resource.

* `storage_class` - (Required, String, ForceNew) Specifies the storage type.  
  The valid values are as follows:
  + **sfs**: SFS file system with v3.0 framework.

  Change this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The NAS storage ID.

* `storage_metadata` - The metadata of the corresponding storage.  
  The [storage_metadata](#workspace_app_nas_storage_metadata) structure is documented below.

* `created_at` - The creation time of the NAS storage.

<a name="workspace_app_nas_storage_metadata"></a>
The `storage_metadata` block supports:

* `export_location` - The storage access URL.

## Import

NAS storages can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_nas_storage.test <id>
```

If the NAS storage ID is unknow, the NAS storage name can be used as an alternative to ID.  

```bash
$ terraform import huaweicloud_workspace_app_nas_storage.test <name>
```
