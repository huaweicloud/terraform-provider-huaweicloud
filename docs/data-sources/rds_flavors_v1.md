---
subcategory: "Deprecated"
---

# huaweicloud\_rds\_flavors\_v1

Use this data source to get the ID of an available HuaweiCloud rds flavor.

!> **Warning:** It has been deprecated, use `huaweicloud_rds_flavors` instead.

## Example Usage

```hcl
data "huaweicloud_rds_flavors_v1" "flavor" {
  region            = "cn-north-1"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.medium"
}
```

## Argument Reference

* `region` - (Required, String) The region in which to obtain the V1 rds client.

* `datastore_name` - (Required, String) The datastore name of the rds.

* `datastore_version` - (Required, String) The datastore version of the rds.

* `speccode` - (Optional, String) The spec code of a rds flavor.

## Available value for attributes

datastore_name | datastore_version | speccode
---- | --- | ---
PostgreSQL | 9.5.5 <br> 9.6.3 | ha = True: <br> rds.pg.m1.2xlarge.ha rds.pg.c2.large.ha rds.pg.s1.2xlarge.ha rds.pg.c2.xlarge.ha rds.pg.s1.xlarge.ha rds.pg.m1.xlarge.ha rds.pg.m1.large.ha rds.pg.c2.medium.ha rds.pg.s1.medium.ha rds.pg.s1.large.ha <br> ha = False: <br> rds.pg.s1.xlarge rds.pg.m1.2xlarge rds.pg.c2.xlarge rds.pg.s1.medium rds.pg.c2.medium rds.pg.s1.large rds.pg.c2.large rds.pg.m1.large rds.pg.s1.2xlarge rds.pg.m1.xlarge
MySQL| 5.6.33 <br>5.6.30  <br>5.6.34 <br>5.6.35 <br>5.6.36 <br>5.7.17| ha = True: <br> rds.mysql.s1.medium.ha rds.mysql.s1.large.ha rds.mysql.s1.xlarge.ha rds.mysql.s1.2xlarge.ha rds.mysql.s1.8xlarge.ha rds.mysql.s1.4xlarge.ha rds.mysql.m1.2xlarge.ha rds.mysql.c2.medium.ha rds.mysql.c2.large.ha rds.mysql.c2.xlarge.ha rds.mysql.c2.2xlarge.ha rds.mysql.c2.4xlarge.ha rds.mysql.c2.8xlarge.ha rds.mysql.m1.medium.ha rds.mysql.m1.large.ha rds.mysql.m1.xlarge.ha rds.mysql.m1.4xlarge.ha <br> ha = False: <br> rds.mysql.s1.medium  rds.mysql.s1.large  rds.mysql.s1.xlarge  rds.mysql.s1.2xlarge  rds.mysql.s1.8xlarge  rds.mysql.s1.4xlarge  rds.mysql.m1.2xlarge  rds.mysql.c2.medium  rds.mysql.c2.large  rds.mysql.c2.xlarge  rds.mysql.c2.2xlarge  rds.mysql.c2.4xlarge  rds.mysql.c2.8xlarge  rds.mysql.m1.medium  rds.mysql.m1.large  rds.mysql.m1.xlarge  rds.mysql.m1.4xlarge
SQLServer| 2014 SP2 SE | <br> ha = True: <br>  rds.mssql.m1.2xlarge.ha rds.mssql.m1.xlarge.ha rds.mssql.m1.4xlarge.ha rds.mssql.s1.xlarge.ha rds.mssql.c2.xlarge.ha rds.mssql.s1.2xlarge.ha <br> ha = False: <br>  rds.mssql.m1.2xlarge  rds.mssql.m1.xlarge  rds.mssql.m1.4xlarge  rds.mssql.s1.xlarge  rds.mssql.c2.xlarge  rds.mssql.s1.2xlarge


## Attributes Reference

`id` is set to the ID of the found rds flavor. In addition, the following attributes
are exported:

* `name` - The name of the rds flavor.
* `ram` - The name of the rds flavor.
