---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_background_tasks"
description: |-
  Use this data source to query the list of DCS background tasks within HuaweiCloud.
---

# huaweicloud_dcs_background_tasks

Use this data source to get the list of DCS background tasks within HuaweiCloud.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dcs_background_tasks" "test" {
  instance_id = var.instance_id
  start_time  = "20240101000000"
  end_time    = "20240131235959"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `start_time` - (Optional, String) Specifies the start time of the query.  
  The time is in UTC format: **yyyyMMddHHmmss**, for example: **20200609160000**.

* `end_time` - (Optional, String) Specifies the end time of the query.  
  The time is in UTC format: **yyyyMMddHHmmss**, for example: **20200609160000**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of background tasks.  
  The [tasks](#dcs_background_tasks_task) structure is documented below.

<a name="dcs_background_tasks_task"></a>
The `tasks` block supports:

* `id` - The background task ID.

* `name` - The background task name.  
  The valid values are as follows:
  + **EXTEND**: Scale specification.
  + **BindEip**: Enable public access.
  + **UnBindEip**: Disable public access.
  + **AddReplica**: Add replica.
  + **DelReplica**: Delete replica.
  + **AddWhitelist**: Set IP whitelist.
  + **UpdatePort**: Modify port.
  + **RemoveIpFromDns**: Remove IP from domain.
  + **masterStandbySwapJob**: Master-standby switchover task.
  + **modify**: Modify password.

* `details` - The detailed information of the task.  
  The [details](#dcs_background_tasks_task_details) structure is documented below.

* `user_name` - The username.

* `user_id` - The user ID.

* `params` - The parameters of the task.

* `status` - The task status.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `error_code` - The error code.

* `enable_show` - Whether task details that can be expanded.

* `job_id` - The task ID.

<a name="dcs_background_tasks_task_details"></a>
The `details` block supports:

* `old_capacity` - The capacity before the modification.

* `new_capacity` - The capacity after the modification.

* `enable_public_ip` - Whether public access has been enabled.

* `public_ip_id` - The public IP address ID.

* `public_ip_address` - The public IP address.

* `enable_ssl` - Whether SSL is enabled.

* `old_cache_mode` - The cache type before the modification.

* `new_cache_mode` - The cache type after the modification.

* `old_resource_spec_code` - The specification parameter before the modification.

* `new_resource_spec_code` - The specification parameter after the modification.

* `old_replica_num` - The number of replicas before the modification.

* `new_replica_num` - The number of replicas after the modification.

* `old_cache_type` - The cache type before the modification.

* `new_cache_type` - The cache type after the modification.

* `replica_ip` - The replica IP address.

* `replica_az` - The AZ where the replica is in.

* `group_name` - The instance shard group name.

* `old_port` - The old port.

* `new_port` - The new port.

* `is_only_adjust_charging` - Whether to only change the billing mode.

* `account_name` - The account name.

* `source_ip` - The source IP address.

* `target_ip` - The target IP address.

* `node_name` - The node information.

* `rename_commands` - The renamed command.

* `updated_config_length` - The length of the updated configuration item.
