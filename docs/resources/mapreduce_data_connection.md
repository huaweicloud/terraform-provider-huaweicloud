---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_data_connection"
description: ""
---

# huaweicloud_mapreduce_data_connection

Manages a data connection resource of Mapreduce within HuaweiCloud.  

## Example Usage

```hcl
variable "db_instance_id" {}
variable "database" {}
variable "user_name" {}
variable "password" {}

resource "huaweicloud_mapreduce_data_connection" "test" {
  name        = "demo"
  source_type = "RDS_MYSQL"
  source_info {
    db_instance_id = var.db_instance_id
    db_name        = var.database
    user_name      = var.user_name
    password       = var.password
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The data connection name.  
  The name can contain a maximum of 64 characters.
  
  Changing this parameter will create a new resource.

* `source_type` - (Required, String, ForceNew) The type of data source.  
  The options are as follows:
    + **RDS_POSTGRES**: RDS PostgreSQL database.
    + **RDS_MYSQL**: RDS MySQL database.
    + **gaussdb-mysql**: GaussDB(for MySQL).
  
  Changing this parameter will create a new resource.

* `source_info` - (Required, List) Information about the data source.  
The [source_info](#DataConnection_SourceInfo) structure is documented below.

<a name="DataConnection_SourceInfo"></a>
The `source_info` block supports:

* `db_instance_id` - (Required, String) The instance ID of database.

* `db_name` - (Required, String) The name of database.

* `user_name` - (Required, String) The user name for logging in to the database.

* `password` - (Required, String) The password for logging in to the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `used_clusters` - Cluster IDs that use this data connection, separated by commas.  

* `status` - The status of the data connection.  
  The valid value are as follows:
    + **0**: data connections are available.

## Import

The data connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_mapreduce_data_connection.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `source_info.0.password`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_mapreduce_data_connection" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_info.0.password,
    ]
  }
}
```
