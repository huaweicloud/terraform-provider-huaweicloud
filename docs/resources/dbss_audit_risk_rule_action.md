---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_audit_risk_rule_action"
description: |-
  Manages a risk rule action resource within HuaweiCloud.
---

# huaweicloud_dbss_audit_risk_rule_action

Manages a risk rule action resource within HuaweiCloud.

-> This resource is only a one-time action resource for doing API action. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "risk_ids" {}
variable "action" {}

resource "huaweicloud_dbss_audit_risk_rule_action" "test" {
  instance_id = var.instance_id
  risk_ids    = var.risk_ids
  action      = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the audit instance ID to which the risk rule belongs.
  Changing this parameter will create a new resource.

* `risk_ids` - (Required, String, ForceNew) Specifies the risk rule ID.
  Multiple IDs should be separated by commas (,).
  Changing this parameter will create a new resource.

* `action` - (Required, String, ForceNew) Specifies the operation type.
  The value can be **ON** or **OFF**.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `result` - The operation execution result.
