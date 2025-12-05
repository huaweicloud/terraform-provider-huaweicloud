---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_operation_connection"
description: |-
  Manages a SecMaster operation connection resource within HuaweiCloud.
---

# huaweicloud_secmaster_operation_connection

Manages a SecMaster operation connection resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}
variable "component_version_id" {}
variable "config" {}

resource "huaweicloud_secmaster_operation_connection" "test" {
  workspace_id         = var.workspace_id
  component_id         = var.component_id
  component_version_id = var.component_version_id
  config               = var.config
  name                 = "test-name"
  description          = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the alert belongs.

* `name` - (Required, String) Specifies the operation connection name.

* `component_id` - (Required, String) Specifies the component ID.

* `component_version_id` - (Required, String) Specifies the component version ID.

* `config` - (Required, String) Specifies the specific operation connection configuration string depends on the
  corresponding field values ​​configured in the plugin's operation connection schema.

* `description` - (Optional, String) Specifies the description of the operation connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the operation connection ID).

* `project_id` - The project ID.

* `component_name` - The component name.

* `type` - The type of assets. Valid values are **datasource** and **action_target**.

* `status` - The status of assets. Valid values are **SUCCESS** and **FAILED**.

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

## Import

The operation connection can be imported using the workspace ID and the operation connection ID, e.g.

```bash
$ terraform import huaweicloud_secmaster_operation_connection.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes includes: `config`.
It is generally recommended running `terraform plan` after importing this resource.
You can then decide if changes should be applied to the resource, or the definition should be updated to align with the
resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_operation_connection" "test" {
  ...

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}
```
