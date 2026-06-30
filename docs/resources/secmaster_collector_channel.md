---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_channel"
description: |-
  Manages a collector channel resource within HuaweiCloud.
---

# huaweicloud_secmaster_collector_channel

Manages a collector channel resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "group_id" {}
variable "parser_id" {}
variable "node_id" {}
variable "input_template_id" {}
variable "input_connection_module_id" {}
variable "output_template_id" {}
variable "output_connection_module_id" {}
variable "input_mode_template_field_id" {}
variable "input_port_template_field_id" {}
variable "input_codec_template_field_id" {}
variable "output_port_template_field_id" {}
variable "output_codec_template_field_id" {}
variable "output_host_template_field_id" {}

resource "huaweicloud_secmaster_collector_channel" "test" {
  workspace_id = var.workspace_id
  title        = "test-channel"
  group_id     = var.group_id
  parser_id    = var.parser_id

  input {
    name                 = "tcp"
    template_id          = var.input_template_id
    connection_module_id = var.input_connection_module_id

    fields {
      name              = "mode"
      type              = "string"
      value             = "server"
      template_field_id = var.input_mode_template_field_id
    }
    fields {
      name              = "port"
      type              = "number"
      value             = "514"
      template_field_id = var.input_port_template_field_id
    }
    fields {
      name              = "codec"
      type              = "string"
      value             = "json"
      template_field_id = var.input_codec_template_field_id
    }
  }

  output {
    name                 = "tcp"
    template_id          = var.output_template_id
    connection_module_id = var.output_connection_module_id

    fields {
      name              = "port"
      type              = "number"
      value             = "514"
      template_field_id = var.output_port_template_field_id
    }
    fields {
      name              = "codec"
      type              = "string"
      value             = "json_lines"
      template_field_id = var.output_codec_template_field_id
    }
    fields {
      name              = "host"
      type              = "string"
      value             = "192.168.0.196"
      template_field_id = var.output_host_template_field_id
    }
  }

  nodes {
    node_id = var.node_id

    args {
      key   = "param1"
      value = "value1"
    }
  }

  description = "created by terraform"
  action      = "INSTALL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `title` - (Required, String) Specifies the collector channel title.

* `group_id` - (Required, String) Specifies the collector channel group ID.

* `parser_id` - (Required, String) Specifies the collector parser ID.

* `input` - (Required, List) Specifies the list of input module configurations.

  The [input](#input_struct) structure is documented below.

* `output` - (Required, List) Specifies the list of output module configurations.

  The [output](#output_struct) structure is documented below.

* `nodes` - (Required, List) Specifies the list of node configurations.

  The [nodes](#nodes_struct) structure is documented below.

* `description` - (Optional, String) Specifies the collector channel description.

* `action` - (Optional, String) Specifies the node operation action.
  The valid values are as follows:
  + **START**: Start.
  + **STOP**: Stop.
  + **REMOVE**: Remove.
  + **RESTART**: Restart.
  + **REFRESH**: Refresh.
  + **INSTALL**: Install.

<a name="input_struct"></a>
The `input` block supports:

* `name` - (Optional, String) Specifies the input module name.

* `template_id` - (Optional, String) Specifies the input module template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated connection module UUID of the input
  module.

* `children` - (Optional, List) Specifies the list of input child modules.

  The [input children](#input_children_struct) structure is documented below.

* `fields` - (Optional, List) Specifies the list of input field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="input_children_struct"></a>
The `input.children` block supports:

* `name` - (Optional, String) Specifies the input child module name.

* `template_id` - (Optional, String) Specifies the input child module template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated connection module UUID of the input
  child module.

* `fields` - (Optional, List) Specifies the list of input child field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="output_struct"></a>
The `output` block supports:

* `name` - (Optional, String) Specifies the output module name.

* `template_id` - (Optional, String) Specifies the output module template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated connection module UUID of the output
  module.

* `children` - (Optional, List) Specifies the list of output child modules.

  The [output children](#output_children_struct) structure is documented below.

* `fields` - (Optional, List) Specifies the list of output field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="output_children_struct"></a>
The `output.children` block supports:

* `name` - (Optional, String) Specifies the output child module name.

* `template_id` - (Optional, String) Specifies the output child module template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated connection module UUID of the output
  child module.

* `fields` - (Optional, List) Specifies the list of output child field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="fields_struct"></a>
The `fields` block supports:

* `name` - (Optional, String) Specifies the field name.

* `type` - (Optional, String) Specifies the field type.

* `value` - (Optional, String) Specifies the field value.

* `other` - (Optional, String) Specifies other supplementary information.

* `template_field_id` - (Optional, String) Specifies the template field UUID.

* `connection_module_id` - (Optional, String) Specifies the associated connection module UUID.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `node_id` - (Optional, String) Specifies the node UUID.

* `node_status` - (Optional, String) Specifies the node status.
  The valid values are **RUN** and **STOP**.

* `args` - (Optional, List) Specifies the list of custom arguments.

  The [args](#args_struct) structure is documented below.

<a name="args_struct"></a>
The `args` block supports:

* `key` - (Optional, String) Specifies the argument key.

* `value` - (Optional, String) Specifies the argument value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the collector channel UUID.

* `create_by` - The IAM user ID.

* `error` - The collector channel error type.

* `operation_status` - The deployment operation status.

* `parser_name` - The collector parser name.

## Import

The collector channel can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_collector_channel.test <workspace_id>/<id>
```
