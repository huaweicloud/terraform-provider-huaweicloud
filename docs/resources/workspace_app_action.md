---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_action"
description: |-
  Manages a Workspace app action resource within HuaweiCloud.
---

# huaweicloud_workspace_app_action

Manages a Workspace app action resource within HuaweiCloud.

-> This resource is used to manage the tenant profile settings. Deleting this resource will not  
   clear the corresponding configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_workspace_app_action" "test" {
  app_restrict_rule_switch = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the tenant profiles are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_restrict_rule_switch` - (Required, Bool, NonUpdatable) Specifies whether to enable the application restriction
  rule switch.
  + **true**: enable
  + **false**: disable

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
