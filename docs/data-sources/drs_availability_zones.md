---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_availability_zones"
description: ""
---

# huaweicloud_drs_availability_zones

Use this data source to query availability zones where DRS jobs can be created within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "migration"
  direction   = "up"
  node_type   = "high"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_type` - (Required, String) Specifies the DRS job engine type.
  Please refer to the document [Engine Types](https://support.huaweicloud.com/intl/en-us/api-drs/drs_api_0159.html).

* `type` - (Required, String) Specifies the job type.

  The options are as follows:
  + **migration**: Online Migration.
  + **sync**: Data Synchronization.
  + **cloudDataGuard**: Disaster Recovery.

* `direction` - (Required, String) Specifies the direction of data flow.

  The options are as follows:
  + **up**: To the cloud. The destination database must be a database in the current cloud.
  + **down**: Out of the cloud. The source database must be a database in the current cloud.
  + **non-dbs**: Self-built database.

* `node_type` - (Required, String) Specifies the node type of the job instance.

  The options are as follows:
  + **micro**: extremely small specification.
  + **small**: small specification.
  + **medium**: medium specification.
  + **high**: large specification.

* `multi_write` - (Optional, Bool) Specifies whether it is dual-AZ disaster recovery.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID.

* `names` - The names of availability zone.
