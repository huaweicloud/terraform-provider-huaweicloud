---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_agency_permissions"
description: |-
  Use this data source to get the list of agency permissions required by DRS service.
---

# huaweicloud_drs_agency_permissions

Use this data source to get the list of agency permissions required by DRS service.

## Example Usage

```hcl
data "huaweicloud_drs_agency_permissions" "test" {
  is_non_dbs = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `is_non_dbs` - (Required, Bool) Specifies whether the database is a self-built database.

* `source_type` - (Optional, String) Specifies the type of the source database.
  The valid values are:
  + **mysql**: MySQL.
  + **sqlserver**: Microsoft SQL Server.
  + **postgresql**: PostgreSQL.
  + **mongodb**: MongoDB, DDS.
  + **oracle**: Oracle.
  + **taurus**: TaurusDB.
  + **ddm**: DDM.
  + **kafka**: Kafka.
  + **gaussdbv5**: GaussDB distributed edition.
  + **gaussdbv5ha**: GaussDB centralized edition.
  + **gaussmongodb**: GeminiDB Mongo.
  + **db2**: DB2.
  + **tidb**: TiDB.
  + **redis**: Redis.
  + **rediscluster**: Redis cluster.
  + **gaussredis**: GeminiDB Redis.
  + **mariadb**: MariaDB.
  + **informix**: Informix.
  + **dynamo**: Dynamo.

* `target_type` - (Optional, String) Specifies the type of the target database.
  The valid values are:
  + **mysql**: MySQL.
  + **sqlserver**: Microsoft SQL Server.
  + **postgresql**: PostgreSQL.
  + **mongodb**: MongoDB, DDS.
  + **oracle**: Oracle.
  + **taurus**: TaurusDB.
  + **ddm**: DDM.
  + **kafka**: Kafka.
  + **gaussdbv5**: GaussDB distributed edition.
  + **gaussdbv5ha**: GaussDB centralized edition.
  + **gaussmongodb**: GeminiDB Mongo.
  + **db2**: DB2.
  + **tidb**: TiDB.
  + **redis**: Redis.
  + **rediscluster**: Redis cluster.
  + **gaussredis**: GeminiDB Redis.
  + **mariadb**: MariaDB.
  + **informix**: Informix.
  + **dynamo**: Dynamo.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `common_permissions` - The list of common permissions required by DRS service.
  The valid values are as follows:
  + **DRS FullAccess**: Full access to Data Replication Service.

* `engine_permissions` - The list of database engine specific permissions.
  The valid values are as follows:
  + **GaussDB ReadOnlyAccess**: Read-only access to GaussDB.
  + **GeminiDB ReadOnlyAccess**: Read-only access to GeminiDB.
  + **DDM ReadOnlyAccess**: Read-only access to DDM.
  + **DDS ReadOnlyPolicy**: Read-only access to DDS.
  + **RDS ReadOnlyAccess**: Read-only access to RDS.
  + **MRS ReadOnlyAccess**: Read-only access to MRS.
