---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_flavors_v3"
sidebar_current: "docs-huaweicloud-datasource-rds-flavors-v3"
description: |-
  Get the flavor information on an HuaweiCloud rds service.
---

# huaweicloud\_rds\_flavors\_v3

Use this data source to get available HuaweiCloud rds flavors.

## Example Usage

```hcl
data "huaweicloud_rds_flavors_v3" "flavor" {
    db_type = "PostgreSQL"
    db_version = "9.5"
    instance_mode = "ha"
}
```

## Argument Reference

* `db_type` - (Required) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `db_version` -
  (Required)
  Specifies the database version. Available value:

type | version
---- | ---
MySQL| 5.6 <br>5.7 <br>8.0
PostgreSQL | 9.5 <br> 9.6 <br>10 <br>11
SQLServer| 2008_R2_EE <br>2008_R2_WEB <br>2012_SE <br>2014_SE <br>2016_SE <br>2017_SE <br>2012_EE <br>2014_EE <br>2016_EE <br>2017_EE <br>2012_WEB <br>2014_WEB <br>2016_WEB <br>2017_WEB

* `instance_mode` - (Required) The mode of instance. Value: ha(indicates primary/standby instance), single(indicates single instance)

## Attributes Reference

In addition, the following attributes are exported:

* `flavors` -
  Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the rds flavor.
* `vcpus` - Indicates the CPU size.
* `memory` - Indicates the memory size in GB.
* `mode` - See 'instance_mode' above.
