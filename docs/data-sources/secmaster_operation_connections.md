---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_operation_connections"
description: |-
  Use this data source to get the list of SecMaster operation connections within HuaweiCloud.
---

# huaweicloud_secmaster_operation_connections

Use this data source to get the list of SecMaster operation connections within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_operation_connections" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the connection name.

* `component_name` - (Optional, String) Specifies the component name.

* `creator_name` - (Optional, String) Specifies the creator name.

* `modifier_name` - (Optional, String) Specifies the modifier name.

* `description` - (Optional, String) Specifies the description.

* `create_start_time` - (Optional, String) Specifies the creation start time.

* `create_end_time` - (Optional, String) Specifies the creation end time.

* `update_start_time` - (Optional, String) Specifies the update start time.

* `update_end_time` - (Optional, String) Specifies the update end time.

* `is_defense_type` - (Optional, Bool) Specifies whether it is used for emergency strategy operation connections.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The operation connections list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The operation connection ID.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `name` - The connection name.

* `component_id` - The component ID.

* `component_name` - The component name.

* `component_version_id` - The component version ID.

* `type` - The type of assets. Valid values are **datasource** and **action_target**.

* `status` - The status of assets. Valid values are **SUCCESS** and **FAILED**.

* `config` - The configuration information.

* `description` - The description.

* `enabled` - Whether the connection is enabled. **false** for disabled, **true** for enabled.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `update_time` - The update time.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `defense_type` - The defense line classification when issuing emergency strategies.

* `target_project_id` - The IAM project ID when issuing emergency strategies.

* `target_project_name` - The IAM project name when issuing emergency strategies.

* `target_enterprise_id` - The enterprise project ID when issuing emergency strategies.

* `target_enterprise_name` - The enterprise project name when issuing emergency strategies.

* `region_id` - The region ID when issuing emergency strategies.

* `region_name` - The region name when issuing emergency strategies.
