---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template"
description: |-
  Manages a DDS parameter template resource within HuaweiCloud.
---

# huaweicloud_dds_parameter_template

Manages a DDS parameter template resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "name" {}
variable "parameter_values" {
  type = map(string)
}
variable "node_type" {}
variable "node_version" {}

resource "huaweicloud_dds_parameter_template" "test"{
  name             = var.name
  parameter_values = var.parameter_values
  node_type        = var.node_type
  node_version     = var.node_version
}
```

### Create a parameter template with entity ID

```hcl
variable "name" {}
variable "entity_id" {}

resource "huaweicloud_dds_parameter_template" "test"{
  name      = var.name
  entity_id = var.entity_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the parameter template name.
  The value must be `1` to `64` characters, which can contain only letters (case-sensitive), digits, hyphens (-),
  underscores (_), and periods (.).

* `node_type` - (Optional, String, ForceNew) Specifies the node type of parameter template.
  Changing this parameter will create a new resource.
  The valid values ara as follows:
  + **mongos**: the mongos node type.
  + **shard**: the shard node type.
  + **config**: the config node type.
  + **replica**: the replica node type.
  + **readonly**: the read replica type of a replica set instance.
  + **shard_readonly**: the read replica type of a cluster instance.
  + **single**: the single node type.

  -> This parameter is mandatory when `entity_id` is not specified.

* `node_version` - (Optional, String, ForceNew) Specifies the database version.
  Changing this parameter will create a new resource.
  The value can be **5.0**, **4.4**, **4.2**, **4.0**, **3.4**.

  -> This parameter is mandatory when `entity_id` is not specified.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can customize parameter values based on the parameters in the default parameter template.

  -> This parameter is mandatory when `entity_id` is not specified.

* `entity_id` - (Optional, String, ForceNew) Specifies the instance ID, group ID, or node ID.

  -> 1.If this parameter is specified, the parameter template is created based on the parameter information of
    the instance, group, or node. The `parameter_values`, `node_type` and `node_version` parameters are ignored.
    <br/>2.If the instance type is cluster, the value is the ID of the shard or config group, ID of the mongos node,
    or ID of the read replica.
    <br/>3.If the instance type is replica set, the value is the instance ID or ID of the read replica.
    <br/>4.If the instance type is single node, the value is the instance ID.

* `description` - (Optional, String) Specifies the parameter template description.
  The description must consist of a maximum of `256` characters and cannot contain the carriage
  return character or the following special characters: >!<"&'=.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `parameters` - Indicates the parameters defined by users based on the default parameter templates.
  The [Parameter](#DdsParameterTemplate_Parameter) structure is documented below.

* `datastore_name` - Indicates database type.

* `created_at` - Indicates the creation time of the parameter template.

* `updated_at` - Indicates the update time of the parameter template.

<a name="DdsParameterTemplate_Parameter"></a>
The `Parameter` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `description` - Indicates the parameter description.

* `type` - Indicates the parameter type.
  + **integer**
  + **string**
  + **boolean**
  + **float**
  + **list**

* `value_range` - Indicates the value range.

* `restart_required` - Indicates whether the instance needs to be restarted.
  + If the value is **true**, restart is required.
  + If the value is **false**, restart is not required.

* `readonly` - Indicates whether the parameter is read-only.
  + If the value is **true**, the parameter is read-only.
  + If the value is **false**, the parameter is not read-only.

## Import

The DDS parameter template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dds_parameter_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `node_type`, `parameter_values`, `entity_id`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dds_parameter_template" "test" {
    ...

  lifecycle {
    ignore_changes = [
      node_type, parameter_values, entity_id,
    ]
  }
}
