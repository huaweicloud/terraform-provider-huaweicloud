---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_memory_rules"
description: |-
  Use this data source to get the list of memory acceleration rules under the specified memory acceleration mapping.
---

# huaweicloud_geminidb_memory_rules

Use this data source to get the list of memory acceleration rules under the specified memory acceleration mapping.

## Example Usage

### Query all rules under the specified mapping

```hcl
variable "mapping_id" {}

data "huaweicloud_geminidb_memory_rules" "test" {
  dbcache_mapping_id = var.mapping_id
}
```

### Query rule by rule ID

```hcl
variable "mapping_id" {}
variable "rule_id" {}

data "huaweicloud_geminidb_memory_rules" "test" {
  dbcache_mapping_id = var.mapping_id
  rule_id            = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `dbcache_mapping_id` - (Required, String) Specifies the memory acceleration mapping ID.

* `rule_id` - (Optional, String) Specifies the memory acceleration rule ID.

* `rule_name` - (Optional, String) Specifies the memory acceleration rule name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

* `source_db_schema` - (Optional, String) Specifies the source database name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

* `source_db_table` - (Optional, String) Specifies the source database table name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of memory acceleration rules.
  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The memory acceleration rule ID.

* `name` - The memory acceleration rule name.

* `status` - The memory acceleration rule status.
  + **normal**
  + **createfail**

* `source_db_schema` - The source database name.

* `source_db_table` - The source database table name.

* `storage_type` - The storage type of the destination database.

* `target_database` - The dstination database.

* `key_columns` - The columns used by mapped keys.

* `value_columns` - The columns used by mapped values.

* `ttl` - The lifetime of a key. (unit: ms)

* `key_separator` - The key separator of a mapping.

* `value_separator` - The value separator of a mapping.

* `key_prefix` - The key prefix.
