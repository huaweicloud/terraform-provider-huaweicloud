---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_high_risk_command"
description: |-
  Manages a GeminiDB high risk command resource within HuaweiCloud.
---

# huaweicloud_geminidb_high_risk_command

Manages a GeminiDB high risk command resource within HuaweiCloud.

-> This resource only supports proxy-based general-purpose GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "huaweicloud_geminidb_high_risk_command" "test" {
  instance_id = var.instance_id
  origin_name = "keys"
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB instance.

* `origin_name` - (Required, String, NonUpdatable) Specifies the original high risk command.
  The valid values are as follows:
  + **keys**
  + **flushdb**
  + **flushall**
  + **hgetall**
  + **hkeys**
  + **hvals**
  + **smembers**

* `name` - (Required, String) Specifies the renamed name of the command.
  If this parameter is left blank, the command is disabled.
  The value can contain a maximum of `30` characters, including digits, letters, and underscores (_).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `instance_id/origin_name`.

## Import

The resource can be imported using the `instance_id` and `origin_name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_geminidb_high_risk_command.test <instance_id>/<origin_name>
```
