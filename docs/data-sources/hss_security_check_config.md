---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_security_check_config"
description: |-
  Use this data source to get the security check configuration of HSS within HuaweiCloud.
---

# huaweicloud_hss_security_check_config

Use this data source to get the security check configuration of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_security_check_config" "test" {}
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

* `task_id` - The scheduled task ID.

* `status` - The security check status. Valid values are:
  + **true**: Enabled.
  + **false**: Disabled.

* `check_period_type` - The check period type. Valid values are:
  + **day**: Daily.
  + **week**: Weekly.

* `day_period` - The daily period.

* `week_period` - The weekly period. Valid values are:
  + **mon**: Monday.
  + **tue**: Tuesday.
  + **wed**: Wednesday.
  + **thu**: Thursday.
  + **fri**: Friday.
  + **sat**: Saturday.
  + **sun**: Sunday.

* `hour` - The check time (hour).

* `content` - The check content. Valid values are:
  + **asset**: Asset.
  + **vul**: Vulnerability.
  + **baseline**: Baseline.
  + **event**: Event.

* `host_id_list` - The list of selected host IDs.
