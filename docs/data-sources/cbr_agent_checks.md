---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_agent_checks"
description: |-
  Use this data source to query the agent status.
---

# huaweicloud_cbr_agent_checks

Use this data source to query the agent status.

## Example Usage

```hcl
variable "resource_id" {}
variable "resource_type" {}

data "huaweicloud_cbr_agent_checks" "test" {
  agent_status {
    resource_id   = var.resource_id
    resource_type = var.resource_type
  }
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the backup detail.
  If omitted, the provider-level region will be used.

* `agent_status` - (Required, List) Specifies the resource Information list.
  The [agent_status](#cbr_agent_status) structure is documented below.

<a name="cbr_agent_status"></a>
The `agent_status` block supports:

* `resource_id` - (Required, String) Specifies the resource ID.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid values are as follows:
  + **OS::Nova::Server**：Indicates ECS.
  + **OS::Ironic::BareMetalServer**：Indicates BMS.

* `resource_name` - (Optional, String) Specifies the resource name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agent_status_attr` - The agent status information list.

  The [agent_status_attr](#agent_status_attr_struct) structure is documented below.

<a name="agent_status_attr_struct"></a>
The `agent_status_attr` block supports:

* `resource_id` - The resource ID.

* `version` - The agent version.

* `installed` - Whether the agent is installed.

* `is_old` - Whether the installed agent is of an earlier version.

* `message` - The error information that explains why the agent cannot be connected.

* `code` - The error code returned uoon an agent connection failure.
