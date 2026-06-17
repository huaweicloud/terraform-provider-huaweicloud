---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_users"
description: |-
  Use this data source to query the database accounts of a TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_users

Use this data source to query the database accounts of a TaurusDB HTAP StarRocks instance within HuaweiCloud.

## Example Usage

### Query all database accounts

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_users" "test" {
  instance_id = var.htap_instance_id
}
```

### Query a specific database account

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_users" "test" {
  instance_id = var.htap_instance_id
  user_name   = "root"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP StarRocks database accounts.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the StarRocks instance ID.

* `user_name` - (Optional, String) Specifies the database account name to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `user_details` - The list of database account details.
  The [user_details](#user_details_attr) structure is documented below.

<a name="user_details_attr"></a>
The `user_details` block supports:

* `user_name` - The database account name.

* `databases` - The list of authorized databases names.

* `dml` - The DML authorization.
  The valid values are as follows:
  + **0**: read and write permissions
  + **1**: read-only permission
  + **2**: read-only and setting permissions
  + **3**: read-write and setting permissions

* `ddl` - The DDL authorization.
  The valid values are as follows:
  + **0**: no DDL permission
  + **1**: DDL permission
