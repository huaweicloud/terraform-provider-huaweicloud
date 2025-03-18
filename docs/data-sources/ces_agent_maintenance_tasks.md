---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_agent_maintenance_tasks"
description: |-
  Use this data source to get the list of CES agent maintenance tasks.
---

# huaweicloud_ces_agent_maintenance_tasks

Use this data source to get the list of CES agent maintenance tasks.

## Example Usage

```hcl
data "huaweicloud_ces_agent_maintenance_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the server ID.

* `instance_type` - (Optional, String) Specifies the server type.
  The valid value can be **ECS** or **BMS**.

* `invocation_id` - (Optional, String) Specifies the task ID.

* `invocation_type` - (Optional, String) Specifies the task type.
  The valid value can be **INSTALL**, **UPDATE**, **ROLLBACK** or **RETRY**.

* `invocation_target` - (Optional, String) Specifies the task object. Only **telescope** is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `invocations` - The task list.

  The [invocations](#invocations_struct) structure is documented below.

<a name="invocations_struct"></a>
The `invocations` block supports:

* `invocation_status` - The task status.

* `create_time` - When the task was created.

* `update_time` - When the task was updated.

* `invocation_id` - The task ID.

* `instance_id` - The server ID.

* `instance_name` - The server name.

* `elastic_ips` - The EIP list.

* `invocation_type` - The task type.

* `current_version` - The current version of the agent.

* `target_version` - The target version.

* `instance_type` - The server type.

* `intranet_ips` - The private IP address list.

* `invocation_target` - The task object.
