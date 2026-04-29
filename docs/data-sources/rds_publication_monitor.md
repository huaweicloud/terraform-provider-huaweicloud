---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publication_monitor"
description: |-
  Use this data source to query the monitor information of an RDS publication within HuaweiCloud.
---

# huaweicloud_rds_publication_monitor

Use this data source to query the monitor information of an RDS publication within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "publication_id" {}

data "huaweicloud_rds_publication_monitor" "test" {
  instance_id    = var.instance_id
  publication_id = var.publication_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the publication monitor.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `publication_id` - (Required, String) Specifies the ID of the publication.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The running status of the snapshot agent associated with the publication.  
  The valid values are as follows:
  + **started**: Started.
  + **succeeded**: Succeeded.
  + **in_progress**: In progress.
  + **idle**: Idle.
  + **retrying**: Retrying.
  + **failed**: Failed.

* `worst_latency` - The longest latency of data changes, in seconds.

* `best_latency` - The shortest latency of data changes, in seconds.

* `average_latency` - The average latency of data changes, in seconds.

* `last_dist_sync` - The last time the distribution agent ran. The format is **yyyy-mm-ddThh:mm:ssZ**.

* `replicated_transactions` - The number of transactions waiting to be transmitted to the distribution database.

* `replication_rate_trans` - The average number of transactions transmitted to the distribution database per second.
