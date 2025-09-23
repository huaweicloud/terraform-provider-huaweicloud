---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_online_data_migration_task"
description: |-
  Manages a DCS online data migration task resource within HuaweiCloud.
---

# huaweicloud_dcs_online_data_migration_task

Manages a DCS online data migration task resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "source_instance_id" {}
variable "target_instance_id" {}

resource "huaweicloud_dcs_online_data_migration_task" "test" {
  task_name         = "test_task_name"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  description       = "terraform test"
  migration_method  = "full_amount_migration"
  resume_mode       = "auto"

  source_instance {
    id       = var.source_instance_id
    password = "test_1234"
  }

  target_instance {
    id       = var.target_instance_id
    password = "test_1234"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `task_name` - (Required, String, NonUpdatable) Specifies the backup import task name.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the network ID of the subnet.

* `security_group_id` - (Required, String, NonUpdatable) Specifies the security group which the instance belongs to.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the online migration task.

* `migration_method` - (Optional, String) Specifies the type of the migration. Value options:
  + **full_amount_migration**: full migration. It is suitable for scenarios where services can be interrupted. Data is
    migrated at one time. Source instance data updated during the migration will not be migrated to the target instance.
  + **incremental_migration**: incremental migration. It is suitable for scenarios requiring minimal service downtime.
    The incremental migration parses logs to ensure data consistency between the source and target instances. After the
    full migration is complete, incremental migration starts.

* `resume_mode` - (Optional, String) Specifies the reconnection mode. Value options:
  + **auto**: automatically reconnect. In this mode, if the source and target instances are disconnected due to network
    exceptions, automatic reconnections will be performed indefinitely. Full synchronization will be triggered and requires
    more bandwidth if incremental synchronization becomes unavailable. Exercise caution when enabling this option.
  + **manual**: manually reconnect.

* `source_instance` - (Optional, List) Specifies the source Redis information.
  The [source_instance](#instance_struct) structure is documented below.

* `target_instance` - (Optional, List) Specifies the target Redis information.
  The [target_instance](#instance_struct) structure is documented below.

* `bandwidth_limit_mb` - (Optional, String) Specifies the bandwidth limit. For incremental migration, you
  can limit the bandwidth to ensure smooth service running. When the data synchronization speed reaches the limit, it
  can no longer increase. Unit: **MB/s**. Value range: **1–10,240** (an integer greater than 0 and less than 10,241).

<a name="instance_struct"></a>
The `source_instance` and `target_instance` block supports:

* `id` - (Optional, String) Specifies the Redis instance ID. It is mandatory if `addrs` is not specified.

* `addrs` - (Optional, String) Specifies the Redis address. It is mandatory if `id` is not specified.

* `password` - (Optional, String) Specifies the Redis password. If a password of the DCS instance is set, it is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `source_instance` - Indicates the source Redis information.
  The [source_instance](#instance_struct) structure is documented below.

* `target_instance` - Indicates the target Redis information.
  The [target_instance](#instance_struct) structure is documented below.

* `ecs_tenant_private_ip` - Indicates the private IP address of the migration ECS on the tenant side. This IP address can
  be added to the whitelist if it is in the same VPC as the private IP address of the target or source Redis.

* `network_type` - Indicates the network type, which can be **VPC** or **VPN**.

* `status` - Indicates the migration task status.

* `supported_features` - Indicates the supported features.

* `version` - Indicates the version of migration ECS.

* `created_at` - Indicates the time when the migration task is created.

* `updated_at` - Indicates the time when the migration task is complete.

* `released_at` - Indicates the time when the migration ECS is released.

<a name="instance_struct"></a>
The `source_instance` and `target_instance` block supports:

* `name` - Indicates the Redis name.

## Timeouts

This resource provides the following timeout configuration option:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DCS backup import task can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dcs_online_data_migration_task.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `source_instance.0.password` and
`target_instance.0.password`. It is generally recommended running `terraform plan` after importing the resource. You can
then decide if changes should be applied to the resource, or the resource definition should be updated to align with the
task. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dcs_online_data_migration_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_instance.0.password, target_instance.0.password,
    ]
  }
}
```
