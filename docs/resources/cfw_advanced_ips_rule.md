---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_advanced_ips_rule"
description: |-
  Manages a resource to create advanced IPS rule within HuaweiCloud.
---

# huaweicloud_cfw_advanced_ips_rule

Manages a resource to create advanced IPS rule within HuaweiCloud.

-> This resource is a one-time action resource used to create advanced IPS rule. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "ips_rule_id" {}
variable "object_id" {}
variable "fw_instance_id" {}

resource "huaweicloud_cfw_advanced_ips_rule" "test" {
  ips_rule_id    = var.ips_rule_id
  object_id      = var.object_id
  fw_instance_id = var.fw_instance_id
  param          = "{\"threshold\":11}"
  action         = 0
  ips_rule_type  = 1
  status         = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, Int, NonUpdatable) Specifies the action type.  
  The valid values are as follows:
  + `0`: Record logs only.
  + `1`: Block session.
  + `2`: Block IP.

* `ips_rule_id` - (Required, String, NonUpdatable) Specifies the advanced IPS rule ID.

* `ips_rule_type` - (Required, Int, NonUpdatable) Specifies the IPS rule type.  
  The valid values are as follows:
  + `0`: Sensitive directory scanning.
  + `1`: Reverse shell xshell.

* `object_id` - (Required, String, NonUpdatable) Specifies the protection object ID.

* `param` - (Required, String, NonUpdatable) Specifies the JSON string containing special parameters.

* `status` - (Required, Int, NonUpdatable) Specifies the rule status.  
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `fw_instance_id` - (Optional, String, NonUpdatable) Specifies the firewall instance ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  For enterprise users, if omitted, default enterprise project will be used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
