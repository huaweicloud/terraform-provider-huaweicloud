---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_shared_folders"
description: |-
  Use this data source to get the list of application shared folders within HuaweiCloud.
---

# huaweicloud_workspace_app_shared_folders

Use this data source to get the list of application shared folders within HuaweiCloud.

## Example Usage

```hcl
variable "storage_id" {}

data "huaweicloud_workspace_app_shared_folders" "test" {
  storage_id = var.storage_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the WKS storage is located.  
  If omitted, the provider-level region will be used.

* `storage_id` - (Required, String) Specifies the WKS storage ID.

* `storage_claim_id` - (Optional, String) Specifies the WKS storage directory claim ID.

* `path` - (Optional, String) Specifies the shared folder path for query.  
  Only visible characters and space(` `) are allowed, must start with a visible character.  
  The valid length is limited from `0` to `128` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `shared_folders` - The list of the shared folders that matched filter parameters.  
  The [shared_folders](#workspace_shared_folders) structure is documented below.

<a name="workspace_shared_folders"></a>
The `shared_folders` block supports:

* `storage_claim_id` - The WKS storage directory claim ID.

* `folder_path` - The storage object path.  

  -> This path is the complete path of the object in the system.

* `delimiter` - The path delimiter.

* `claim_mode` - The storage claim type.
  + **USER**: user directory
  + **SHARE**: shared directory.

* `count` - The number of associated users and user groups of the shared folder.  
  It is a map with keys like `USER` and `USER_GROUP`, and values are integers representing the counts.
