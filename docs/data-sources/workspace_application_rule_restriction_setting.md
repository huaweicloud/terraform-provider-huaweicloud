---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rule_restriction_setting"
description: |-
  Use this data source to get the application rule restriction setting within HuaweiCloud.
---

# huaweicloud_workspace_application_rule_restriction_setting

Use this data source to get the application rule restriction setting within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_application_rule_restriction_setting" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application rule restriction setting is located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `app_restrict_rule_switch` - Whether the application restriction rule switch is enabled.

* `app_control_mode` - The application control mode.

* `app_periodic_switch` - Whether the periodic monitoring switch is enabled.

* `app_periodic_interval` - The periodic monitoring interval time, in minutes.

* `app_force_kill_proc_switch` - Whether the force kill application switch is enabled.
