---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_databases"
description: |-
  Use this data source to query the list of databases of a TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_databases

Use this data source to query the list of databases of a TaurusDB HTAP StarRocks instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_databases" "test" {
  instance_id = var.htap_instance_id
}
```

### Filter by database name

```hcl
variable "htap_instance_id" {}
variable "htap_database_name" {}

data "huaweicloud_taurusdb_htap_starrocks_databases" "test" {
  instance_id   = var.htap_instance_id
  database_name = var.htap_database_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP StarRocks databases.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP StarRocks instance ID.

* `database_name` - (Optional, String) Specifies the database name to be queried. The name supports fuzzy query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The list of database names.
