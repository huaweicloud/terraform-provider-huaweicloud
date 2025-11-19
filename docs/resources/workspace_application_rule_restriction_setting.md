---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rule_restriction_setting"
description: |-
  Use this resource to setting the application rule restriction within HuaweiCloud.
---

# huaweicloud_workspace_application_rule_restriction_setting

Use this resource to setting the application rule restriction within HuaweiCloud.

-> This resource is only a one-time action resource for setting the application rule restriction. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Set application rule restriction with full parameters

```hcl
resource "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  app_restrict_rule_switch   = true
  app_control_mode           = 0
  app_periodic_switch        = true
  app_periodic_interval      = 10
  app_force_kill_proc_switch = true
}
```

### Set application rule restriction with minimum parameters

```hcl
resource "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  app_restrict_rule_switch = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application rule restriction is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_restrict_rule_switch` - (Required, Bool, NonUpdatable) Specifies whether to enable the application restriction
  rule switch.

* `app_control_mode` - (Optional, Int, NonUpdatable) Specifies the application control mode.  
  The value `0` indicates that applications in the list are prohibited from running.

* `app_periodic_switch` - (Optional, Bool, NonUpdatable) Specifies whether to enable the periodic monitoring switch.

* `app_periodic_interval` - (Optional, Int, NonUpdatable) Specifies the periodic monitoring interval time, in minutes.  
  The minmum value is `5`.

* `app_force_kill_proc_switch` - (Optional, Bool, NonUpdatable) Specifies whether to enable the force kill application
  switch.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
