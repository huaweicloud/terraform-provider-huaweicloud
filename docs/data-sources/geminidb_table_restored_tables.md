---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_table_restored_tables"
description: |-
  Use this data source to obtain GeminiDB Cassandra instance table information that is restored using tables.
---

# huaweicloud_geminidb_table_restored_tables

Use this data source to obtain GeminiDB Cassandra instance table information that is restored using tables.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_table_restored_tables" "test" {
  instance_id   = var.instance_id
  database_name = "test_db"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the tables.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GeminiDB Cassandra instance.

* `database_name` - (Required, String) Specifies the name of the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_names` - The list of table names.
