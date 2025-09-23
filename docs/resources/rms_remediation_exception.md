---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_remediation_exception"
description: |-
  Manages a RMS remediation exception resource within HuaweiCloud.
---

# huaweicloud_rms_remediation_exception

Manages a RMS remediation exception resource within HuaweiCloud.

## Example Usage

```hcl
variable "policy_assignment_id" {}
variable "resource_id" {}

resource "huaweicloud_rms_remediation_exception" "test" {
  policy_assignment_id = var.policy_assignment_id
  
  exceptions {
    resource_id = var.resource_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `policy_assignment_id` - (Required, String, NonUpdatable) Specifies the policy assignment ID.

* `exceptions` - (Required, List) Specifies the list of remediation exceptions.
  The [exceptions](#Exceptions) structure is documented below.

<a name="Exceptions"></a>
The `exceptions` block supports:

* `resource_id` - (Required, String) Specifies the resource ID.

* `message` - (Optional, String) Specifies the reason for adding an exception.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the same as the `policy_assignment_id`.

* `exceptions` - The list of remediation exceptions.
  The [exceptions](#ExceptionsAttr) structure is documented below.

<a name="ExceptionsAttr"></a>
The `exceptions` block supports:

* `created_by` - The creator of a remediation exception.

* `joined_at` - The time when a remediation exception is added.

## Import

The remediation exception can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_rms_remediation_exception.test <id>
```
