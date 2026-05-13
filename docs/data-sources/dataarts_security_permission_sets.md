---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_permission_sets"
description: |-
  Use this data source to get the list of permission sets for DataArts Studio Security within HuaweiCloud.
---

# huaweicloud_dataarts_security_permission_sets

Use this data source to get the list of permission sets for DataArts Studio Security within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_permission_sets" "test" {
  workspace_id = var.workspace_id
}
```

### Query permission sets by parent ID

```hcl
variable "workspace_id" {}
variable "parent_id" {}

data "huaweicloud_dataarts_security_permission_sets" "test" {
  workspace_id = var.workspace_id
  parent_id    = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the permission sets are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the permission sets belong.

* `name` - (Optional, String) Specifies the name of the permission set to be queried.  
  This parameter supports fuzzy matching.

* `parent_id` - (Optional, String) Specifies the parent ID of the permission set to be queried.

* `type_filter` - (Optional, String) Specifies the type filter of the permission sets to be queried.  
  The valid values are as follows:
  + **TOP_PERMISSION_SET**
  + **SUB_PERMISSION_SET**
  + **ALL_PERMISSION_SET**

* `manager_id` - (Optional, String) Specifies the manager ID of the permission set to be queried.

* `manager_name` - (Optional, String) Specifies the manager name of the permission set to be queried.

* `manager_type` - (Optional, String) Specifies the manager type of the permission set to be queried.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**

* `datasource_type` - (Optional, String) Specifies the datasource type of the permission set to be queried.  
  The valid values are as follows:
  + **HIVE**
  + **DWS**
  + **DLI**

* `sync_status` - (Optional, String) Specifies the sync status of the permission set to be queried.  
  The valid values are as follows:
  + **UNKNOWN**
  + **NOT_SYNC**
  + **SYNCING**
  + **SYNC_SUCCESS**
  + **SYNC_FAIL**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permission_sets` - The list of permission sets that match the filter parameters.  
  The [permission_sets](#dataarts_security_permission_sets_attr) structure is documented below.

<a name="dataarts_security_permission_sets_attr"></a>
The `permission_sets` block supports:

* `id` - The ID of the permission set.

* `parent_id` - The parent ID of the permission set.

* `name` - The name of the permission set.

* `description` - The description of the permission set.

* `type` - The type of the permission set.  
  + **COMMON**
  + **MRS_MANAGED**

* `managed_cluster_id` - The managed cluster ID of the permission set.

* `managed_cluster_name` - The managed cluster name of the permission set.

* `project_id` - The project ID of the permission set.

* `domain_id` - The domain ID of the permission set.

* `instance_id` - The instance ID of the permission set.

* `manager_id` - The manager ID of the permission set.

* `manager_name` - The manager name of the permission set.

* `manager_type` - The manager type of the permission set.

* `datasource_type` - The datasource type of the permission set.

* `sync_status` - The sync status of the permission set.

* `sync_msg` - The sync message of the permission set.

* `sync_time` - The sync time of the permission set, in RFC3339 format.

* `created_at` - The creation time of the permission set, in RFC3339 format.

* `updated_at` - The update time of the permission set, in RFC3339 format.

* `created_by` - The creator of the permission set.

* `updated_by` - The updater of the permission set.
