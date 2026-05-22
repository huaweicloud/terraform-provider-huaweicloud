---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_primary_instance_tables"
description: |-
  Use this data source to query the list of tables of a TaurusDB HTAP primary instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_primary_instance_tables

Use this data source to query the list of tables of a TaurusDB HTAP primary instance within HuaweiCloud.

## Example Usage

### Filter by database and table with whitelist

```hcl
variable "htap_instance_id" {}
variable "source_taurusdb_instance_id" {}

data "huaweicloud_taurusdb_htap_primary_instance_tables" "test" {
  instance_id        = var.htap_instance_id
  source_instance_id = var.source_taurusdb_instance_id
  filter_type        = "include_tables"

  database_tables {
    database = "db"
    tables   = ["table1"]
  }

  selected_tables {
    database = "db"
    tables   = ["table1", "table2"]
  }
}
```

### Filter by database and table with blacklist

```hcl
variable "htap_instance_id" {}
variable "source_taurusdb_instance_id" {}

data "huaweicloud_taurusdb_htap_primary_instance_tables" "test" {
  instance_id        = var.htap_instance_id
  source_instance_id = var.source_taurusdb_instance_id
  filter_type        = "exclude_tables"

  database_tables {
    database = "db"
    tables   = ["table3"]
  }

  selected_tables {
    database = "db"
    tables   = ["table1", "table2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP primary instance tables.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP instance ID.

* `source_instance_id` - (Required, String) Specifies the ID of the primary instance whose tables are to be queried.

* `database_tables` - (Required, List) Specifies the list of database and table names of the primary instance to be retrieved.
  The [database_tables](#database_tables_attr) structure is documented below.

* `selected_tables` - (Required, List) Specifies the list of databases and tables names to be selected.
  The [selected_tables](#database_tables_attr) structure is documented below.

* `filter_type` - (Required, String) Specifies the table filter type, blacklist or whitelist.
  The valid values are as follows:
  + **include_tables**: whitelist.
  + **exclude_tables**: blacklist.

<a name="database_tables_attr"></a>
The `database_tables` block supports:

* `database` - (Required, String) Specifies the name of a database to be queried.

* `tables` - (Required, List) Specifies the names of the data tables to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - The list of table names. The name in format of `<database_anme>.<table_name>`.
