---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_database_versions"
description: |-
  Use this data source to query the GeminiDB database version information within HuaweiCloud.
---

# huaweicloud_geminidb_database_versions

Use this data source to query the GeminiDB database version information within HuaweiCloud.

## Example Usage

```hcl
variable "datastore_name" {}

data "huaweicloud_geminidb_database_versions" "test" {
  datastore_name = var.datastore_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `datastore_name` - (Required, String) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: Indicates GeminiDB Cassandra.
  + **mongodb**: Indicates GeminiDB Mongo.
  + **influxdb**: Indicates GeminiDB Influx.
  + **redis**: Indicates GeminiDB Redis.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The list of database versions.
