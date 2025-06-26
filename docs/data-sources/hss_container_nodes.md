---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_nodes"
description: |-
  Use this data source to get the list of HSS container nodes within HuaweiCloud.
---

# huaweicloud_hss_container_nodes

Use this data source to get the list of HSS container nodes within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_nodes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS container nodes.
  If omitted, the provider-level region will be used.

* `host_name` - (Optional, String) Specifies the host name.

* `agent_status` - (Optional, String) Specifies the agent status.  
  The valid values are as follows:
  + **not_installed**
  + **online**
  + **offline**

* `protect_status` - (Optional, String) Specifies the protection status.  
  The valid values are as follows:
  + **closed**
  + **opened**

* `container_tags` - (Optional, String) Specifies the label used to identify CCE container nodes or self built nodes.  
  The valid values are as follows:
  + **cce**: CCE nodes.
  + **self**: Self built nodes.
  + **other**: Other nodes.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - All container nodes that match the filter parameters.  
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_status` - The host status.  
  The valid values are as follows:
  + **ACTIVE**
  + **SHUTOFF**
  + **BUILDING**
  + **ERROR**

* `agent_status` - The agent status.  
  The valid values are as follows:
  + **not_installed**
  + **online**
  + **offline**

* `protect_status` - The protection status.
  The valid values are as follows:
  + **closed**
  + **opened**

* `protect_interrupt` - Is the protection interrupted.

* `protect_degradation` - Has the protection been downgraded.

* `degradation_reason` - The reasons for downgraded protection.

* `container_tags` - The label used to identify CCE container nodes or self built nodes.  
  The valid values are as follows:
  + **cce**: CCE node.
  + **self**: Self built nodes.
  + **other**: Other nodes.

* `private_ip` - The private ip address.

* `public_ip` - The elastic public IP address.

* `resource_id` - The host security quota ID (UUID).

* `group_name` - The host group name.

* `enterprise_project_name` - The enterprise project name.

* `detect_result` - The host security detection results.  
  The valid values are as follows:
  + **undetected**: Not detected.
  + **clean**: No Risk.
  + **risk**: At risk.
  + **scanning**: Detecting.

* `asset` - The asset risk.

* `vulnerability` - The vulnerability risk.

* `intrusion` - The intrusion risk.

* `policy_group_id` - The policy group ID.

* `policy_group_name` - The policy group name.
