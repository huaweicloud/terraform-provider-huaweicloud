---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_overview_agent_statistics"
description: |-
  Use this data source to get the agent statistics of HSS within HuaweiCloud.
---

# huaweicloud_hss_overview_agent_statistics

Use this data source to get the agent statistics of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_overview_agent_statistics" "test" {}
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

* `container_type` - (Optional, Int) Specifies whether it is a container asset.
  The value can be:
  + **0**: Not a container asset.
  + **1**: Container asset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `wait_upgrade_num` - The number of agents waiting for upgrade.

* `online_num` - The number of online agents.

* `not_online_num` - The number of not online agents.

* `offline_num` - The number of offline agents.

* `incluster_num` - The number of nodes in the cluster.

* `not_incluster_num` - The number of nodes not in the cluster.
