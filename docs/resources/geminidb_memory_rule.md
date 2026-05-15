---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_memory_rule"
description: |-
  Manages a memory acceleration rule resource within HuaweiCloud.
---

# huaweicloud_geminidb_memory_rule

Manages a memory acceleration rule resource within HuaweiCloud.

-> This resource only supports primary/standby GeminiDB Redis instance.

## Example Usage

```hcl
var "mapping_id" {}
var "rule_name" {}
var "database_name" {}
var "table_name" {}
var "key_prefix" {}
var "key_columns" {
  type = list(string)
}
var "value_columns" {
  type = list(string)
}

resource "huaweicloud_geminidb_memory_rule" "test" {
  dbcache_mapping_id = var.mapping_id
  name               = var.rule_name
  source_db_schema   = var.database_name
  source_db_table    = table_name
  storage_type       = "hash"
  target_database    = "0"
  key_prefix         = var.key_prefix
  key_columns        = var.key_columns
  value_columns      = var.value_columns
  key_separator      = "."
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `dbcache_mapping_id` - (Required, String, NonUpdatable) Specifies the ID of the memory acceleration mapping.

* `name` - (Required, String, NonUpdatable) Specifies the name of the memory acceleration rule.
  The value can contain a maximum of `256` characters and must be unique in the current mapping.

* `source_db_schema` - (Required, String, NonUpdatable) Specifies the source database.

* `source_db_table` - (Required, String, NonUpdatable) Specifies the source database table.

* `storage_type` - (Required, String, NonUpdatable) Specifies the storage type of the destination database.
  The value only can be **hash**.

* `target_database` - (Required, String, NonUpdatable) Specifies the destination database.
  The valid value rangs from `0` to `999`.

* `key_prefix` - (Required, String, NonUpdatable) Specifies the key prefix.
  The value can contain no more than `1,024` characters.
  
-> The parameters `key_prefix` and `target_database` must be unique in the current mapping.

* `key_columns` - (Required, List, NonUpdatable) Specifies the columns used by mapped keys.

* `value_columns` - (Required, List) Specifies the columns used by mapped values.

* `key_separator` - (Required, String, NonUpdatable) Specifies the key separator of a mapping.
  Only one character is allowed.

* `value_separator` - (Optional, String) Specifies the value separator of a mapping.
  Only one character is allowed.

* `ttl` - (Optional, String) Specifies the lifetime of a key, in ms.
  If this parameter is not specified, default value `2,592,000,000` is used, indicating 30 days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The memory acceleration rule status.
  + **normal**
  + **createfail**

## Import

The resource can be imported using the `dbcache_mapping_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_geminidb_memory_rule.test <dbcache_mapping_id>/<id>
```
