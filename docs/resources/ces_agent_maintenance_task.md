---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_agent_maintenance_task"
description: |-
  Manages a CES agent maintenance task within HuaweiCloud.
---

# huaweicloud_ces_agent_maintenance_task

Manages a CES agent maintenance task within HuaweiCloud.

## Example Usage

### Create an install task

```hcl
variable "instance_id" {}

resource "huaweicloud_ces_agent_maintenance_task" "test" {
  invocation_type = "INSTALL"
  instance_id     = var.instance_id
}
```

### Create an update task

```hcl
variable "instance_id" {}

resource "huaweicloud_ces_agent_maintenance_task" "test" {
  invocation_type = "UPDATE"
  instance_id     = var.instance_id
  version_type    = "ADVANCE_VERSION"
  version         = "2.7.5.1"
}
```

### Create a rollback task

```hcl
variable "invocation_id" {}

resource "huaweicloud_ces_agent_maintenance_task" "test" {
  invocation_type = "ROLLBACK"
  invocation_id   = var.invocation_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `invocation_type` - (Required, String, NonUpdatable) Specifies the task type.
  The valid value can be **INSTALL**, **UPDATE**, **ROLLBACK** or **RETRY**.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the server ID.
  This parameter is mandatory when the task type is **INSTALL** or **UPDATE**.

* `invocation_target` - (Optional, String, NonUpdatable) Specifies the task object. Only **telescope** is supported.

* `invocation_id` - (Optional, String, NonUpdatable) Specifies the task ID.
  This parameter is mandatory when the task type is **ROLLBACK** or **RETRY**.

* `version_type` - (Optional, String, NonUpdatable) Specifies the version the agent will be upgraded to.
  The valid value can be **BASIC_VERSION** or **ADVANCE_VERSION**.

* `version` - (Optional, String, NonUpdatable) Specifies the version number.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `invocations` - The result of the agent maintenance task creation.
  The [invocations](#Invocations) structure is documented below.

<a name="Invocations"></a>
The `invocations` block supports:

* `invocation_id` - The task ID.

* `instance_id` - The server ID.

* `instance_name` - The server name

* `instance_type` - The server type.

* `intranet_ips` - The private IP address list.

* `elastic_ips` - The EIP list.

* `invocation_type` - The task type.

* `invocation_status` - The task status.

* `invocation_target` - The task object.

* `create_time` - When the task was created.

* `update_time` - When the task was updated.

* `current_version` - The current version of the agent.

* `target_version` - The target version.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
