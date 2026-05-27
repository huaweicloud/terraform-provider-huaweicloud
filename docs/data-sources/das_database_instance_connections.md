---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_connections"
description: |-
  Use this data source to query DAS Database instance connections within HuaweiCloud.
---

# huaweicloud_das_instance_connections

Use this data source to query DAS Database instance connections within HuaweiCloud.

## Example Usage

### Query all database instance connections

```hcl
data "huaweicloud_das_instance_connections" "test" {}
```

### Query database instance connections using keyword search

```hcl
# Fuzzy search by database username, name, address, or remarks, which target fields contain this keyword.
variable "fuzzy_search_keyword" {}

data "huaweicloud_das_instance_connections" "by_keyword" {
  condition = var.fuzzy_search_keyword 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database instance connections are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance to which the database instance connection
  belongs.

* `network_type` - (Optional, String) Specifies the network type of the database instance connection.  
  The valid values are as follows:
  + **rds**
  + **gaussdb**
  + **dds**
  + **ddm**

* `datastore_type` - (Optional, String) Specifies the datastore type of the database instance connection.  
  The valid values are as follows:
  + **mysql**
  + **sqlserver**
  + **postgresql**
  + **taurus**
  + **gaussdbv5**
  + **mongodb**
  + **ddm**

* `connection_type` - (Optional, String) Specifies the connection type of the database instance connection.  
  The valid values are as follows:
  + **NORMAL**
  + **SHARE**

* `condition` - (Optional, String) Specifies the keyword used to search for database instance connection address, name,
  database username, or remarks.  
  The fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `connections` - The list of connections that matched filter parameters.  
  The [connections](#das_database_instance_connections) structure is documented below.

<a name="das_database_instance_connections"></a>
The `connections` block supports:

* `id` - The ID of the database instance connection.

* `instance_id` - The instance ID of the database instance connection.

* `engine_type` - The engine type of the database instance connection.

* `network_type` - The network type of the database instance connection.

* `username` - The username of the database instance connection.

* `is_save_password` - Whether to save the password for the database instance connection.

* `description` - The description of the database instance connection.

* `port` - The port of the database instance connection.

* `database_name` - The database name of the database instance connection.

* `instance_name` - The instance name of the database instance connection.

* `datastore_version` - The datastore version of the database instance connection.

* `ip_address` - The ip address of the database instance connection.

* `created_at` - The time when the database instance connection was created, in RFC3339 format.

* `status` - The status of the database instance connection.

* `conn_share_type` - The conn share type of the database instance connection.

* `shared_user_name` - The shared user name of the database instance connection.

* `shared_user_id` - The shared user ID of the database instance connection.

* `expired_at` - The time when the database instance connection expires, in RFC3339 format.
