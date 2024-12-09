---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_nas_storages"
description: |-
  Use this data source to get the list of NAS storages for Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_nas_storages

Use this data source to get the list of NAS storages for Workspace APP within HuaweiCloud.

## Example Usage

### Query all NAS storages

```hcl
data "huaweicloud_workspace_app_nas_storages" "test" {}
```

### Query NAS storage containing fragments with the same name

```hcl
variable "storage_name_prefix" {}

data "huaweicloud_workspace_app_nas_storages" "test" {
  name = var.storage_name_prefix
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the NAS storages are located.  
  If omitted, the provider-level region will be used.

* `storage_id` - (Optional, String) Specifies the ID of the NAS storage to be queried.

* `name` - (Optional, String) Specifies the name of the NAS storage to be queried.  
  This parameter is used for fuzzy search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storages` - All NAS storages that match the filter parameters.
  The [storages](#workspace_app_nas_storages) structure is documented below.

<a name="workspace_app_nas_storages"></a>
The `storages` block supports:

* `id` - The ID of the NAS storage.

* `name` - The name of the NAS storage.

* `storage_metadata` - The metadata of the corresponding storage.
  The [storage_metadata](#workspace_app_nas_storage_metadata_attr) structure is documented below.

* `created_at` - The creation time of the NAS storage, in RFC3339 format.

* `personal_folder_count` - The number of the personal folders under this NAS storage.

* `shared_folder_count` - The number of the shared folders under this NAS storage.

<a name="workspace_app_nas_storage_metadata_attr"></a>
The `storage_metadata` block supports:

* `storage_handle` - The storage name.

* `storage_class` - The storage type.

* `export_location` - The storage access URL.
