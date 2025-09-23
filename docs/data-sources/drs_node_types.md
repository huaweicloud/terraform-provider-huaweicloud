---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_node_types"
description: |-
  Use this data source to query available node types for DRS jobs.
---

# huaweicloud_drs_node_types

Use this data source to query available node types for DRS jobs.

## Example Usage

```hcl
data "huaweicloud_drs_node_types" "test" {
  engine_type = "mysql"
  type        = "sync"
  direction   = "up"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `engine_type` - (Required, String) Specifies the DRS job engine type.
  For details, see [engine types](https://support.huaweicloud.com/intl/en-us/api-drs/drs_api_0159.html)

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

* `multi_write` - (Optional, Bool) Specifies whether dual-active disaster recovery is used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_types` - The available node types list.
