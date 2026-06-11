---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_command_disable"
description: |-
  Manages a GeminiDB Redis command disabled resource within HuaweiCloud.
---

# huaweicloud_geminidb_command_disable

Manages a memory acceleration rule resource within HuaweiCloud.

-> This resource only supports GeminiDB Redis instance.

## Example Usage

### create a disable command resource

```hcl
var "instance_id" {}
var "commands" {
  type = list(string)
}

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = var.instance_id
  disabled_type = "command"
  commands      = var.commands
}
```

### create a disable key command resource

```hcl
var "instance_id" {}
variable "keys" {
  type = list(object({
    db_id    = string
    key      = string
    commands = list(string)
  }))
}

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = var.instance_id
  disabled_type = "key"
  
  dynamic "keys" {
    for_each = var.keys

    content {
      db_id    = keys.value.db_id
      key      = keys.value.key
      commands = keys.value.command
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB Redis instance.

* `disabled_type` - (Required, String, NonUpdatable) Specifies the disable type.
  The valid values are as follows:
  + **command**
  + **key**

* `commands` - (Optional, List) Specifies the list of disabled commands.
  The valid values are as follows:
  + **keys**
  + **hkeys**
  + **hvals**
  + **hgetall**
  + **smembers**
  + **flushdb**
  + **flushall**

  -> This parameter is available and required when the `disabled_type` is set to **command**.

* `keys` - (Optional, List) Specifies the list of keys information. A maximum of `20` keys are supported.
  The [keys](#keys_struct) structure is documented below.

  -> This parameter is available and required when the `disabled_type` is set to **key**.

<a name="keys_struct"></a>
The `keys` block supports:

* `db_id` - (Required, String) Specifies the DB ID where a key is located.

* `key` - (Required, String) Specifies the key name.

* `commands` - (Required, List) Specifies the list of commands.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the instance ID.

## Import

The resource can be imported using the `instance_id` and `disabled_type`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_geminidb_command_disable.test <instance_id>/<disabled_type>
```
