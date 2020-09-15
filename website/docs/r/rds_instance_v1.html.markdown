---
layout: "huaweicloud"
page_title: "OpenStack: huaweicloud_rds_instance_v1"
sidebar_current: "docs-huaweicloud-resource-rds-instance-v1"
description: |-
  Manages rds instance resource within HuaweiCloud
---

# huaweicloud\_rds\_instance\_v1

Manages rds instance resource within HuaweiCloud

## Example Usage:  Creating a PostgreSQL RDS instance

```hcl
data "huaweicloud_rds_flavors" "flavor" {
  region            = "eu-de"
  datastore_name    = "PostgreSQL"
  datastore_version = "9.5.5"
  speccode          = "rds.pg.s1.large"
}

resource "huaweicloud_networking_secgroup" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "huaweicloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = data.huaweicloud_rds_flavors.flavor.id
  volume {
    type = "COMMON"
    size = 200
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = "c1095fe7-03df-4205-ad2d-6f4c181d436e"
  nics {
    subnetid = "b65f8d25-c533-47e2-8601-cfaa265a3e3e"
  }
  securitygroup {
    id = huaweicloud_networking_secgroup.secgrp_rds.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "04:00:00"
    keepdays  = 4
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable          = true
    replicationmode = "async"
  }
  depends_on = ["huaweicloud_networking_secgroup.secgrp_rds"]
}
```

## Example Usage:  Creating a SQLServer RDS instance
```hcl
data "huaweicloud_rds_flavors" "flavor" {
  region            = "eu-de"
  datastore_name    = "SQLServer"
  datastore_version = "2014 SP2 SE"
  speccode          = "rds.mssql.s1.2xlarge"
}

resource "huaweicloud_networking_secgroup" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "huaweicloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "SQLServer"
    version = "2014 SP2 SE"
  }
  flavorref = data.huaweicloud_rds_flavors.flavor.id
  volume {
    type = "COMMON"
    size = 200
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = "c1095fe7-03df-4205-ad2d-6f4c181d436e"
  nics {
    subnetid = "b65f8d25-c533-47e2-8601-cfaa265a3e3e"
  }
  securitygroup {
    id = huaweicloud_networking_secgroup.secgrp_rds.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "04:00:00"
    keepdays  = 4
  }
  dbrtpd     = "Huangwei!120521"
  depends_on = ["huaweicloud_networking_secgroup.secgrp_rds"]
}
```

## Example Usage:  Creating a MySQL RDS instance
```hcl
data "huaweicloud_rds_flavors" "flavor" {
  region            = "eu-de"
  datastore_name    = "MySQL"
  datastore_version = "5.6.33"
  speccode          = "rds.mysql.s1.medium"
}

resource "huaweicloud_networking_secgroup" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "huaweicloud_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type    = "MySQL"
    version = "5.6.33"
  }
  flavorref = data.huaweicloud_rds_flavors.flavor.id
  volume {
    type = "COMMON"
    size = 200
  }
  region           = "eu-de"
  availabilityzone = "eu-de-01"
  vpc              = "c1095fe7-03df-4205-ad2d-6f4c181d436e"
  nics {
    subnetid = "b65f8d25-c533-47e2-8601-cfaa265a3e3e"
  }
  securitygroup {
    id = huaweicloud_networking_secgroup.secgrp_rds.id
  }
  dbport = "8635"
  backupstrategy {
    starttime = "04:00:00"
    keepdays  = 4
  }
  dbrtpd = "Huangwei!120521"
  ha {
    enable          = true
    replicationmode = "async"
  }
  depends_on = ["huaweicloud_networking_secgroup.secgrp_rds"]
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the DB instance name. The DB instance name of
    the same type is unique in the same tenant.

* `datastore` - (Required) Specifies database information. The structure is
    described below.

* `flavorref` - (Required) Specifies the specification ID (flavors.id in the
    response message in Obtaining All DB Instance Specifications).

* `volume` - (Required) Specifies the volume information. The structure is described
    below.

* `region` - (Required) Specifies the region ID.

* `availabilityzone` - (Required) Specifies the ID of the AZ.

* `vpc` - (Required) Specifies the VPC ID. For details about how to obtain this
    parameter value, see section "Virtual Private Cloud" in the Virtual Private
    Cloud API Reference.

* `nics` - (Required) Specifies the nics information. For details about how
    to obtain this parameter value, see section "Subnet" in the Virtual Private
    Cloud API Reference. The structure is described below.

* `securitygroup` - (Required) Specifies the security group which the RDS DB
    instance belongs to. The structure is described below.

* `dbport` - (Optional) Specifies the database port number.

* `backupstrategy` - (Optional) Specifies the advanced backup policy. The structure
    is described below.

* `dbrtpd` - (Required) Specifies the password for user root of the database.

* `ha` - (Optional) Specifies the parameters configured on HA and is used when
    creating HA DB instances. The structure is described below. NOTICE:
    RDS for Microsoft SQL Server does not support creating HA DB instances and
    this parameter is not involved.

The `datastore` block supports:

* `type` - (Required) Specifies the DB engine. Currently, MySQL, PostgreSQL, and
    Microsoft SQL Server are supported. The value is MySQL, PostgreSQL, or SQLServer.

* `version` - (Required) Specifies the DB instance version.


Available value for attributes

type | version
---- | ---
PostgreSQL | 9.5.5 <br> 9.6.3
MySQL| 5.6.33 <br>5.6.30  <br>5.6.34 <br>5.6.35 <br>5.6.36 <br>5.7.17
SQLServer| 2014 SP2 SE


The `volume` block supports:

* `type` - (Required) Specifies the volume type. Valid value:
    It must be COMMON (SATA) or ULTRAHIGH (SSD) and is case-sensitive.

* `size` - (Required) Specifies the volume size.
    Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB.

The `nics` block supports:

* `subnetId` - (Required) Specifies the subnet ID obtained from the VPC.

The `securitygroup ` block supports:

* `id` - (Required) Specifies the ID obtained from the securitygroup.

The `backupstrategy ` block supports:

* `starttime` - (Optional) Indicates the backup start time that has been set.
    The backup task will be triggered within one hour after the backup start time.
    Valid value: The value cannot be empty. It must use the hh:mm:ss format and
    must be valid. The current time is the UTC time.

* `keepdays` - (Optional) Specifies the number of days to retain the generated backup files.
    Its value range is 0 to 35. If this parameter is not specified or set to 0, the
    automated backup policy is disabled.

The `ha` block supports:

* `enable` - (Optional) Specifies the configured parameters on the HA.
    Valid value: The value is true or false. The value true indicates creating
    HA DB instances. The value false indicates creating a single DB instance.

* `replicationmode` - (Optional) Specifies the replication mode for the standby DB instance.
    The value cannot be empty.
    For MySQL, the value is async or semisync.
    For PostgreSQL, the value is async or sync.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `flavorref` - See Argument Reference above.
* `volume` - See Argument Reference above.
* `availabilityzone` - See Argument Reference above.
* `vpc` - See Argument Reference above.
* `nics` - See Argument Reference above.
* `securitygroup` - See Argument Reference above.
* `dbport` - See Argument Reference above.
* `backupstrategy` - See Argument Reference above.
* `dbrtpd` - See Argument Reference above.
* `ha` - See Argument Reference above.
* `status` - Indicates the DB instance status.
* `hostname` - Indicates the instance connection address. It is a blank string.
* `type` - Indicates the DB instance type, which can be master or readreplica.
* `created` - Indicates the creation time in the following format: yyyy-mm-dd Thh:mm:ssZ.
* `updated` - Indicates the update time in the following format: yyyy-mm-dd Thh:mm:ssZ.

## Attributes Reference

The following attributes can be updated:

* `volume.size` - See Argument Reference above.

* `flavorref` - See Argument Reference above.

* `backupstrategy` - See Argument Reference above.
