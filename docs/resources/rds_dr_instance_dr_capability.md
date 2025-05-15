---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_dr_instance_dr_capability"
description: |-
  Manages RDS dr instance dr capability resource within HuaweiCloud.
---

# huaweicloud_rds_dr_instance_dr_capability

Manages RDS dr instance dr capability resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "target_instance_id" {}
variable "target_project_id" {}
variable "target_region" {}
variable "target_ip" {}

resource "huaweicloud_rds_dr_instance_dr_capability" "test" {
  instance_id        = var.instance_id
  target_instance_id = var.target_instance_id
  target_project_id  = var.target_project_id
  target_region      = var.target_region
  target_ip          = var.target_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `target_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the primary DB instance.

* `target_project_id` - (Required, String, NonUpdatable) Specifies the project ID of the tenant to which the primary DB
  instance belongs.

* `target_region` - (Required, String, NonUpdatable) Specifies the ID of the region where the primary DB instance resides.

* `target_ip` - (Required, String, NonUpdatable) Specifies the data VIP of the primary DB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the DR configuration status.
* `replica_state` - Indicates the Synchronization status. The value can be:
  + **0**: indicates that the synchronization is normal.
  + **-1**: indicates that the synchronization is abnormal.

* `time` - Indicates the DR configuration time.

* `build_process` - Indicates the process for configuring disaster recovery (DR). The value can be:
  + **master**: process of configuring DR capability for the primary instance
  + **slave**: process of configuring DR for the DR instance

* `wal_receive_replay_delay_in_ms` - Indicates the replay delay, in milliseconds, on the DR instance.

* `wal_write_receive_delay_in_mb` - Indicates the WAL send lag volume, in MB. It means the difference between the WAL Log
  Sequence Number (LSN) written by the primary instance and the WAL LSN received by the DR instance.

* `wal_write_replay_delay_in_mb` - Indicates the end-to-end delayed WAL size, in MB. It refers to the difference between
  the WAL LSN written by the primary instance and the WAL LSN replayed by the DR instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The RDS dr instance dr capability can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_dr_instance_dr_capability.test <id>
```
