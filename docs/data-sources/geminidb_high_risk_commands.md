---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_high_risk_commands"
description: |-
  Use this data source to get the list of high risk commmands.
---

# huaweicloud_geminidb_high_risk_commands

Use this data source to get the list of high risk commmands.

-> This data source only supports proxy-based general-purpose GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_high_risk_commands" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `commands` - The list of high risk commands.
  The [commands](#commands_struct) structure is documented below.

<a name="commands_struct"></a>
The `commands` block supports:

* `origin_name` - The original high risk command name.

* `name` - The renamed name of the command.
