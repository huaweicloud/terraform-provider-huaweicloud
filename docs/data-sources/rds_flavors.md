---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_flavors"
description: |-
  Use this data source to get available rds flavors.
---

# huaweicloud_rds_flavors

Use this data source to get available rds flavors.

## Example Usage

```hcl
data "huaweicloud_rds_flavors" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "9.5"
  instance_mode = "ha"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `db_type` - (Required, String) Specifies the DB engine. The value can be **MySQL**, **PostgreSQL**, **SQLServer**,
  **MariaDB**.

* `db_version` - (Optional, String) Specifies the database version.

* `instance_mode` - (Optional, String) Specifies the mode of instance. Value options:
  + **ha**: indicates primary/standby instance
  + **single**: indicates single instance
  + **replica**: indicates read replicas

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the RDS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the RDS flavor.

* `group_type` - (Optional, String) Specifies the performance specification, the valid values are as follows:
  + **normal**: General enhanced.
  + **normal2**: General enhanced type II.
  + **armFlavors**: KunPeng general enhancement.
  + **dedicatedNormal**: (dedicatedNormalLocalssd): Dedicated for x86.
  + **armLocalssd**: KunPeng general type.
  + **normalLocalssd**: x86 general type.
  + **general**: General type.
  + **dedicated**:  
    For MySQL engine: Dedicated type.  
    For PostgreSQL and SQL Server engines: Dedicated type, only supported by cloud disk SSD.
  + **rapid**:  
    For MySQL engine: Dedicated (discontinued).  
    For PostgreSQL and SQL Server engines: Dedicated, only supported by ultra-fast SSDs.
  + **bigmem**: Large memory type.
  + **yunyao**: Flexus RDS type.

* `availability_zone` - (Optional, String) Specifies the availability zone which the RDS flavor belongs to.

* `is_flexus` - (Optional, Bool) Specifies whether to query flexus RDS instance specifications.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the list of flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ID of the rds flavor.

* `name` - Indicates the name of the rds flavor.

* `vcpus` - Indicates the CPU size.

* `memory` - Indicates the memory size in GB.

* `group_type` - Indicates the performance specification.

* `instance_mode` - Indicates the mode of instance.

* `availability_zones` - Indicates the availability zones which the RDS flavor belongs to.

* `db_versions` - Indicates the Available versions of the database.

* `az_status` - Indicates the specification status in an AZ.~~
