---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_statistic"
description: |-
  Use this data source to get the HSS baseline statistic within HuaweiCloud.
---

# huaweicloud_hss_baseline_statistic

Use this data source to get the HSS baseline statistic within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_statistic" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_id` - (Optional, String) Specifies the server ID.

* `group_id` - (Optional, String) Specifies the policy group ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `host_weak_pwd` - The number of weak password detections.

* `pwd_policy` - The number of password complexity strategy detections.

* `security_check` - The number of server configuration checks.
