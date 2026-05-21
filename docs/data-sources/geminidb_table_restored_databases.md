---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_table_restored_databases"
description: |-
  Use this data source to obtain GeminiDB Cassandra instance database information that is restored using tables.
---

# huaweicloud_geminidb_table_restored_databases

Use this data source to obtain GeminiDB Cassandra instance database information that is restored using tables.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_table_restored_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the databases.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GeminiDB Cassandra instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_names` - The list of database names.
