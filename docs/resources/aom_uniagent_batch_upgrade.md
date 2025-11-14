---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_uniagent_batch_upgrade"
description: |-
  Use this resource to batch upgrade AOM UniAgent within HuaweiCloud.
---

# huaweicloud_aom_uniagent_batch_upgrade

Use this resource to batch upgrade AOM UniAgent within HuaweiCloud.

-> This resource is only a one-time action resource for batch upgrading AOM UniAgents. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
variable "version" {}
variable "agent_list" {
  type = list(object({
    agent_id = string
    inner_ip = string
  }))
}

resource "huaweicloud_aom_uniagent_batch_upgrade" "test" {
  version = var.version

  dynamic "agent_list" {
    for_each = var.agent_list

    content {
      agent_id = agent_list.value.agent_id
      inner_ip = agent_list.value.inner_ip
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the target machines to be operated are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `version` - (Required, String, NonUpdatable) Specifies the version number of UniAgent to be upgraded.

* `agent_list` - (Required, List, NonUpdatable) Specifies the list of host information for upgrading UniAgent.  
  Up to `100` hosts are supported.  
  The [agent_list](#aom_agent_list) structure is documented below.

<a name="aom_agent_list"></a>
The `agent_list` block supports:

* `agent_id` - (Required, String) Specifies the unique agent ID.

* `inner_ip` - (Required, String) Specifies the host IP address.
