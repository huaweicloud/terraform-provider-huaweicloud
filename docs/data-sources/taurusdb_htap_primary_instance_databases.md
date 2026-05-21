---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_primary_instance_databases"
description: |-
  Use this data source to query the databases of a primary HTAP instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_primary_instance_databases

Use this data source to query the databases of a primary HTAP instance within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}
variable "source_taurusdb_instance_id" {}
variable "taurusdb_databases" {
  type = list(string)
}

data "huaweicloud_taurusdb_htap_primary_instance_databases" "test" {
  instance_id         = var.htap_instance_id
  source_instance_id  = var.source_taurusdb_instance_id
  databases           = var.taurusdb_databases
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the databases.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP instance ID.

* `databases` - (Required, List) Specifies the list of database names of the primary instance to be retrieved.
  The name in list supports fuzzy query.

* `source_instance_id` - (Required, String) Specifies the ID of the primary instance whose database is to be retrieved.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_names` - The database names.
