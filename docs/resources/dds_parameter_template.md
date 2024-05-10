---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template"
description: ""
---

# huaweicloud_dds_parameter_template

Manages a DDS parameter template resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "parameter_values" {}
variable "node_type" {}
variable "node_version" {}

resource "huaweicloud_dds_parameter_template" "test"{
  name             = var.name
  parameter_values = var.parameter_values
  node_type        = var.node_type
  node_version     = var.node_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the parameter template name.
  The value must be 1 to 64 characters, which can contain only letters, digits, hyphens (-),
  underscores (_), and periods (.).

* `node_type` - (Required, String, ForceNew) Specifies the node type of parameter template. Valid value:
  + **mongos**: the mongos node type.
  + **shard**: the shard node type.
  + **config**: the config node type.
  + **replica**: the replica node type.
  + **single**: the single node type.

  Changing this parameter will create a new resource.

* `node_version` - (Required, String, ForceNew) Specifies the database version.
  The value can be **4.4**, **4.2**, **4.0**, **3.4** or **3.2**.

  Changing this parameter will create a new resource.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can customize parameter values based on the parameters in the default parameter template.

* `description` - (Optional, String) Specifies the parameter template description.
  The description must consist of a maximum of 256 characters and cannot contain the carriage
  return character or the following special characters: >!<"&'=.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `parameters` - Indicates the parameters defined by users based on the default parameter templates.
  The [Parameter](#DdsParameterTemplate_Parameter) structure is documented below.

* `created_at` - The create time of the parameter template.

* `updated_at` - The update time of the parameter template.

<a name="DdsParameterTemplate_Parameter"></a>
The `Parameter` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `description` - Indicates the parameter description.

* `type` - Indicates the parameter type. The value can be integer, string, boolean, float, or list.

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
$ terraform import huaweicloud_dds_parameter_template.test <tempalate_id>
```
