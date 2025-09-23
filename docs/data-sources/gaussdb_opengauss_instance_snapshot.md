---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_snapshot"
description: |-
  Use this data source to get the information About the Original Instance Based on a Specific Point of Time or a Backup File.
---

# huaweicloud_gaussdb_opengauss_instance_snapshot

Use this data source to get the information About the Original Instance Based on a Specific Point of Time or a Backup File.

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_gaussdb_opengauss_instance_snapshot" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the original instance ID.
  It is mandatory if `restore_time` is specified.

* `restore_time` - (Optional, String) Specifies the instance information at a time point.
  It is in the UNIX timestamp format, in milliseconds. The time zone is UTC. It is mandatory when you want to view DB
  instance backups based on a specified point in time.

* `backup_id` - (Optional, String) Specifies the backup ID.
  It is mandatory when a DB instance is restored using a backup ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cluster_mode` - Indicates the instance deployment model.
  The value can be:
  + **Ha**: primary/standby deployment
  + **Independent**: independent deployment
  + **Combined**: combined deployment

* `instance_mode` - Indicates the instance model.
  The value can be:
  + **basic**: basic edition
  + **standard**: standard edition
  + **enterprise**: enterprise edition

* `data_volume_size` - Indicates the storage space, in GB

* `solution` - Indicates the solution template type.
  The value can be:
  + **single**: single node
  + **double**: 1 primary + 1 standby (2 nodes)
  + **triset**: 1 primary + 2 standby
  + **logger**: 1 primary + 1 standby + 1 log
  + **loggerdorado**: 1 primary + 1 standby + 1 log (shared storage)
  + **quadruset**: 1 primary + 3 standby
  + **hws**: distributed (independent deployment)

* `node_num` - Indicates the number of nodes.

* `coordinator_num` - Indicates the number of CNs.

* `sharding_num` - Indicates the number of shards.

* `replica_num` - Indicates the number of replicas.

* `engine_version` - Indicates the engine version.
