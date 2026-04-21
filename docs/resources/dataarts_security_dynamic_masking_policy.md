---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_dynamic_masking_policy"
description: |-
  Manages a dynamic masking policy resource for DataArts Studio Security within HuaweiCloud.
---

# huaweicloud_dataarts_security_dynamic_masking_policy

Manages a dynamic masking policy resource for DataArts Studio Security within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "policy_name" {}
variable "dws_cluster_id" {}
variable "dws_cluster_name" {}
variable "database_name" {}
variable "table_name" {}
variable "connection_name" {}
variable "connection_id" {}
variable "users" {}
variable "user_groups" {}
variable "schema_name" {}

resource "huaweicloud_dataarts_security_dynamic_masking_policy" "test" {
  workspace_id    = var.workspace_id
  name            = var.policy_name
  datasource_type = "DWS"
  cluster_id      = var.dws_cluster_id
  cluster_name    = var.dws_cluster_name
  database_name   = var.database_name
  table_name      = var.table_name
  conn_name       = var.connection_name
  conn_id         = var.connection_id
  users           = var.users
  user_groups     = var.user_groups
  schema_name     = var.schema_name

  policy_list {
    column_name          = "name"
    column_type          = "text"
    algorithm_type       = "DWS_SELF_CONFIG"
    algorithm_detail_dto = jsonencode({
      start         = 1
      end           = 2
      string_target = "*"
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dynamic masking policy is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the dynamic masking
  policy belongs.

* `name` - (Required, String) Specifies the name of the dynamic masking policy.  
  The valid length is limited from `2` to `64`, only letters, Chinese characters, digits, and underscores (_) are allowed,
  and must start with letters or Chinese characters.

* `datasource_type` - (Required, String, NonUpdatable) Specifies the data source type of the dynamic masking policy.  
  The valid values are as follows:
  + **HIVE**
  + **DWS**
  + **DLI**

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster corresponding to the data source.  
  When the data source type is **DLI**, this parameter must be set to **DLI**.

* `cluster_name` - (Required, String, NonUpdatable) Specifies the name of the cluster corresponding to the
  data source.  
  When the data source type is **DLI**, this parameter must be set to **DLI**.

* `database_name` - (Required, String, NonUpdatable) Specifies the name of the database.

* `table_name` - (Required, String, NonUpdatable) Specifies the name of the data table.

* `conn_name` - (Required, String) Specifies the name of the data connection.

* `conn_id` - (Required, String) Specifies the ID of the data connection.

* `policy_list` - (Required, List) Specifies the list of dynamic masking policy configurations.  
  The [policy_list](#dataarts_security_dynamic_masking_policy_list) structure is documented below.

* `table_id` - (Optional, String) Specifies the ID of the data table.

* `user_groups` - (Optional, String) Specifies the list of user groups, separated by commas (,).  
  At least one of the `user_groups` and `users` parameters must be specified.

* `users` - (Optional, String) Specifies the list of users, separated by commas (,).  
  At least one of the `user_groups` and `users` parameters must be specified.
  
* `schema_name` - (Optional, String, NonUpdatable) Specifies the schema name corresponding to the DWS data source.

<a name="dataarts_security_dynamic_masking_policy_list"></a>
The `policy_list` block supports:

* `column_name` - (Required, String) Specifies the field name in the data table.

* `column_type` - (Required, String) Specifies the field type in the data table.

* `algorithm_type` - (Optional, String) Specifies the algorithm type of dynamic masking.  
  For HIVE data source dynamic masking algorithm, the valid values are as follows:
  + **MASK**: Mask letters and digits.
  + **MASK_SHOW_LAST_4**: Show last 4 characters.
  + **MASK_SHOW_FIRST_4**: Show first 4 characters.
  + **MASK_HASH**: Hash masking.
  + **MASK_DATE_SHOW_YEAR**: Mask month and date.
  + **MASK_NULL**: Null masking.

  For DWS data source dynamic masking algorithm, the valid values are as follows:
  + **DWS_ALL_MASK**: Full masking.
  + **DWS_BACK_KEEP**: Keep last 4 characters, mask the rest as *.
  + **DWS_FRONT_KEEP**: Keep first 2 characters, mask the rest as *.
  + **DWS_SELF_CONFIG**: Custom masking. The start and end positions and the masking character need to be specified.

  For DLI data source dynamic masking algorithm, the valid values are as follows:
  + **MASK**: Mask letters and digits.
  + **MASK_SHOW_LAST_4**: Keep last 4 characters.
  + **MASK_SHOW_FIRST_4**: Keep first 4 characters.
  + **MASK_HASH**: Hash masking.
  + **MASK_DATE_SHOW_YEAR**: Mask month and date.
  + **MASK_NULL**: Null masking.

* `algorithm_detail_dto` - (Optional, String) Specifies the algorithm detail object of dynamic masking, in
  JSON format.  
  For field details, please refer to the [documentation](https://support.huaweicloud.com/api-dataartsstudio/CreateSecurityDynamicMaskingPolicy.html#CreateSecurityDynamicMaskingPolicy__request_AlgorithmDetailDTO).
  
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `sync_status` - The current synchronization status of the policy.
  + **UNKNOWN**
  + **NOT_SYNC**
  + **SYNC_SUCCESS**
  + **SYNC_FAIL**
  + **SYNC_PARTIAL_FAIL**
  + **DATA_UPDATED**

* `sync_msg` - The synchronization message of the policy.

* `sync_log` - The synchronization log of the policy.

* `create_time` - The creation time of the policy, in RFC3339 format.

* `create_user` - The creator of the policy.

* `update_time` - The latest update time of the policy, in RFC3339 format.

* `update_user` - The latest updater of the policy.

## Import

The resource can be imported using `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_dynamic_masking_policy.test <workspace_id>/<id>
```
