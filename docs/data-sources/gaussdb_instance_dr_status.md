---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_dr_status"
description: |-
  Use this data source to query the disaster recovery status of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_dr_status

Use this data source to query the disaster recovery status of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_instance_dr_status" "test" {
  instance_id   = var.instance_id
  disaster_type = "stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID to query the DR status.

* `disaster_type` - (Required, String) Specifies the disaster recovery type.
  The value can be **stream** (streaming disaster recovery).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.

* `status` - The disaster recovery status.
  The valid values are:
  + **normal**: Normal
  + **creating**: Creating
  + **create_fail**: Creation failed
  + **data_sync**: Data synchronizing
  + **stopping**: Stopping
  + **stopped**: Stopped
  + **stop_fail**: Stop failed
  + **recovering**: Recovering
  + **recovery_fail**: Recovery failed
  + **upgrading**: Upgrading
  + **promote**: Promoting to primary
  + **promote_fail**: Promote to primary failed

* `rpo` - The Recovery Point Objective (RPO) value.

* `rto` - The Recovery Time Objective (RTO) value.

* `rpo_threshold` - The RPO threshold.

* `rto_threshold` - The RTO threshold.

* `switchover_progress` - The switchover progress (percentage).

* `failover_progress` - The failover progress (percentage).
