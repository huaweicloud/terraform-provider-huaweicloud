---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_agent_statuses"
description: |-
  Use this data source to get the list of agent statuses.
---

# huaweicloud_ces_agent_statuses

Use this data source to get the list of agent statuses.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_ces_agent_statuses" "test" {
  instance_ids = [var.instance_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_ids` - (Required, List) Specifies the cloud server ID list.

* `uniagent_status` - (Optional, String) Specifies the uniagent status.
  The valid value can be **none** (not installed), **running**, **silent** or **unknown** (faulty).

* `extension_name` - (Optional, String) Specifies the agent name.
  If this parameter is not specified, all agents are queried.
  Currently, only telescope can be queried.

* `extension_status` - (Optional, String) Specifies the agent status.
  If this parameter is not specified, all statuses are queried.
  The valid value can be **none** (not installed), **running**, **stopped**, **fault** (process exception)
  or **unknown** (connection exception).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agent_status` - The agent statuses.

  The [agent_status](#agent_status_struct) structure is documented below.

<a name="agent_status_struct"></a>
The `agent_status` block supports:

* `instance_id` - The cloud server ID.

* `uniagent_status` - The uniagent status.

* `extensions` - The agent extension information list.

  The [extensions](#agent_status_extensions_struct) structure is documented below.

<a name="agent_status_extensions_struct"></a>
The `extensions` block supports:

* `name` - The agent name.

* `status` - The agent status.

* `version` - The agent version.
