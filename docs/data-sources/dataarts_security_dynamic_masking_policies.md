---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_dynamic_masking_policies"
description: |-
  Use this data source to get the list of DataArts Security dynamic masking policies within HuaweiCloud.
---

# huaweicloud_dataarts_security_dynamic_masking_policies

Use this data source to get the list of DataArts Security dynamic masking policies within HuaweiCloud.

## Example Usage

### Query all dynamic masking policies under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_dynamic_masking_policies" "test" {
  workspace_id = var.workspace_id
}
```

### Query dynamic masking policies by name

```hcl
variable "workspace_id" {}
variable "policy_name" {}

data "huaweicloud_dataarts_security_dynamic_masking_policies" "test" {
  workspace_id = var.workspace_id
  name         = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the dynamic masking policies.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the dynamic masking policies belong.

* `name` - (Optional, String) Specifies the name of the dynamic masking policy to be queried.  
  Fuzzy search is supported.

* `cluster_name` - (Optional, String) Specifies the name of the cluster to be queried.  
  Fuzzy search is supported.

* `database_name` - (Optional, String) Specifies the name of the database to be queried.  
  Fuzzy search is supported.

* `table_name` - (Optional, String) Specifies the name of the data table to be queried.  
  Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of dynamic masking policies that matched filter parameters.  
  The [policies](#dataarts_security_dynamic_masking_policies_attr) structure is documented below.

<a name="dataarts_security_dynamic_masking_policies_attr"></a>
The `policies` block supports:

* `id` - The ID of the dynamic masking policy.

* `name` - The name of the dynamic masking policy.

* `datasource_type` - The data source type of the dynamic masking policy.
  + **DWS**
  + **DLI**
  + **HIVE**

* `cluster_id` - The ID of the cluster corresponding to the data source.

* `cluster_name` - The name of the cluster corresponding to the data source.

* `database_name` - The name of the database.

* `table_name` - The name of the data table.

* `user_groups` - The user groups of the dynamic masking policy, separated by commas (,).

* `users` - The users of the dynamic masking policy, separated by commas (,).

* `sync_status` - The synchronization status of the dynamic masking policy.  
  + **UNKNOWN**
  + **NOT_SYNC**
  + **SYNC_SUCCESS**
  + **SYNC_FAIL**
  + **SYNC_PARTIAL_FAIL**
  + **DATA_UPDATED**

* `sync_time` - The synchronization time of the dynamic masking policy, in RFC3339 format.

* `sync_msg` - The synchronization log of the dynamic masking policy.

* `create_time` - The creation time of the dynamic masking policy, in RFC3339 format.

* `create_user` - The creator of the dynamic masking policy.

* `update_time` - The latest update time of the dynamic masking policy, in RFC3339 format.

* `update_user` - The latest updater of the dynamic masking policy.
