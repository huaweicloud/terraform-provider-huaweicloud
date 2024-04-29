---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_nosql_flavors"
description: ""
---

# huaweicloud_gaussdb_nosql_flavors

Use this data source to get available HuaweiCloud GeminiDB flavors.
This is an alternative to `huaweicloud_gaussdb_cassandra_flavors`

## Example Usage

```hcl
data "huaweicloud_gaussdb_nosql_flavors" "flavors" {
  vcpus  = 4
  memory = 8
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the GaussDB specifications.
  If omitted, the provider-level region will be used.

* `engine` - (Optional, String) Specifies the type of the database engine. The valid values are as follows:
  + **cassandra**: The default value and means to query GaussDB (for Cassandra) instance specifications.
  + **redis**: Means to query GaussDB (for Redis) instance specifications.
  + **mongodb**: Means to query GaussDB (for Mongo) instance specifications.
  + **influxdb**: Means to query GaussDB (for Influx) instance specifications.

* `engine_version` - (Optional, String) Specifies the version of the database engine.

* `vcpus` - (Optional, Int) Specifies the number of vCPUs.

* `memory` - (Optional, Int) Specifies the memory size in gigabytes (GB).

* `availability_zone` - (Optional, String) Specifies the availability zone (AZ) of the GaussDB specifications.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `flavors` - The information of the GaussDB specifications. Structure is documented below.

The `flavors` block contains:

* `name` - The spec code of the flavor.

* `vcpus` - The number of vCPUs.

* `memory` - The memory size, in GB.

* `engine` - The type of the database engine.

* `engine_version` - The version of the database engine.

* `availability_zones` - All available zones (on sale) for current flavor.
