---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rule_restriction_switch"
description: |-
  Use this resource to enable or disable the application rule restriction within HuaweiCloud.
---

# huaweicloud_workspace_application_rule_restriction_switch

Use this resource to enable or disable the application rule restriction within HuaweiCloud.

-> This resource is only a one-time action resource for enabling or disabling the application rule restriction. Deleting
   this resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Enable Application Rule Restriction

```hcl
resource "huaweicloud_workspace_application_rule_restriction_switch" "test" {
  action = "enable"
}
```

### Disable Application Rule Restriction

```hcl
resource "huaweicloud_workspace_application_rule_restriction_switch" "test" {
  action = "disable"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application rule restriction to be operated is
  located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the action type for the application rule restriction.  
  Valid values are:
  + **enable**: Enable the application rule restriction.
  + **disable**: Disable the application rule restriction.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
