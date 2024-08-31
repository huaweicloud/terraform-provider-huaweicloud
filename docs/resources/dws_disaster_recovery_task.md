---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_disaster_recovery_task"
description: |-
  Manages a GaussDB(DWS) disaster recovery task resource within HuaweiCloud.
---

# huaweicloud_dws_disaster_recovery_task

Manages a GaussDB(DWS) disaster recovery task resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "primary_cluster_id" {}
variable "standby_cluster_id" {}
variable "sync_period" {}

resource "huaweicloud_dws_disaster_recovery_task" "test" {
  name               = var.name
  dr_type            = "az"
  primary_cluster_id = var.primary_cluster_id
  standby_cluster_id = var.standby_cluster_id
  dr_sync_period     = var.sync_period
  action             = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `dr_type` - (Required, String, ForceNew) Specifies the type of the DR task. Only support **az** now.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DR task. It must be unique and
  contains `4` to `64` characters, which consist of letters, digits, hyphens(-), or underscores(_) only
  and must start with a letter. Changing this creates a new resource.

* `primary_cluster_id` - (Required, String, ForceNew) Specifies the ID of the primary cluster.
  Changing this creates a new resource.

* `standby_cluster_id` - (Required, String, ForceNew) Specifies the ID of the standby cluster.
  Changing this creates a new resource.

* `dr_sync_period` - (Required, String) Specifies the synchronization period of the DR task. The valid
  value ranges from `1` to `3000`, the uint support `m` minute, `H` hour,`d` day. e.g. **20m**, means 20 minutes.
  When `status` is **unstart** or **stopped**, `dr_sync_period` can be change.

* `action` - (Optional, String) Specifies the action for the DR task. The valid values are:
  + **start**: Starting the DR task. You can start a DR task when `status` is **unstart**, **stopped**  
   or **start_failed**.
  + **pause**: Stopping the DR task. You can stop a DR task when `status` is **running** or **stop_failed**.
  + **switchover**: Switching to the DR cluster. You can perform a DR switchover when `status` is **running**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time, in UTC format.

* `primary_cluster` - The primary cluster.
  The [cluster](#attrblock_cluster) structure is documented below.

* `standby_cluster` - The standby cluster.
  The [cluster](#attrblock_cluster) structure is documented below.

* `started_at` - The start time of the DR task, in UTC format.

* `status` - The status of the DR task.

<a name="attrblock_cluster"></a>
The `cluster` block supports:

* `id` - The cluster ID.

* `name` - The cluster name.

* `cluster_az` - The availability zone to which the cluster belongs.

* `last_success_at` - The lasted success synchronized time, in UTC format.

* `obs_bucket_name` - The cluster OBS name.

* `progress` - The DR task cluster progress.

* `region` - The region to which the cluster belongs.

* `role` - The DR task cluster role.

* `status` - The DR task cluster status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DR task can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dws_disaster_recovery.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `action`, `primary_cluster_id`, `standby_cluster_id`.
It is generally recommended running `terraform plan` after importing a DR task.
You can then decide if changes should be applied to the DR task, or the resource definition
should be updated to align with the DR task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_disaster_recovery" "test" {
    ...

  lifecycle {
    ignore_changes = [
      action, primary_cluster_id, standby_cluster_id
    ]
  }
}
```
