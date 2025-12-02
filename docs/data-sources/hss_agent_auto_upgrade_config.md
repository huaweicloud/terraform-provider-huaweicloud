---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_agent_auto_upgrade_config"
description: |-
  Use this data source to get the agent auto upgrade configuration of HSS within HuaweiCloud.
---

# huaweicloud_hss_agent_auto_upgrade_config

Use this data source to get the agent auto upgrade configuration of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_agent_auto_upgrade_config" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `enabled` - The agent auto upgrade configuration switch status.  
  The valid values are as follows:
  + **true**: Enabled.
  + **false**: Disabled.
