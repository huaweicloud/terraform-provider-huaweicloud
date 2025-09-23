---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_data_connection"
description: ""
---

# huaweicloud_dataarts_studio_data_connection

Using this resource you can configure data source information and create a data connection within HuaweiCloud,
through which you can access data sources when developing scripts and jobs.

-> There can be only one data connection of the same type.

## Example Usage

### Create a data connection of the DLI type

```hcl
variable "workspace_id" {}
variable "connection_name" {}

resource "huaweicloud_dataarts_studio_data_connection" "test" {
  workspace_id = var.workspace_id
  type         = "DLI"
  name         = var.connection_name
  env_type     = 0
  config       = jsonencode({
    "cdm_property_enable": "false"
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the data connection is located.
  Changing this will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the data connection belongs.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the data connection name.  
  The value contains a maximum of `100` characters, including only letters, digits, hyphens (-), and underscores (_).  
  The connection name must be unique.

* `type` - (Required, String) Specifies the data connection type.

* `config` - (Optional, String) Specifies the dynamic configuration for the specified type of data connection.

  -> Please fill the dynamic configuration based on the debugging result on the console. Pay attention to the type of
     the type of the value: the `bool` values (HCL format) and the `bool values of string type` (HCL format) are not
     considered equal.

* `agent_id` - (Optional, String) Specifies the agent ID.

* `agent_name` - (Optional, String) Specifies the agent name.  
  Required if the `agent_id` is specified.

* `env_type` - (Optional, Int) Specifies the data connection mode.
  + **0**: Development mode.
  + **1**: Production mode.

  Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `qualified_name` - The qualified name of this data connection.

* `created_at` - The creation time of the data connection.

* `created_by` - The name of the data connection creator.

## Import

Data connections can be imported using related `workspace_id` and their `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_studio_data_connection.test <workspace_id>/<name>
```
