---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_batch_delete_template_rule"
description: |-
  Manages a resource to batch delete template rule associations within HuaweiCloud.
---

# huaweicloud_dsc_batch_delete_template_rule

Manages a resource to batch delete template rule associations within HuaweiCloud.

-> This resource is a one-time action resource used to batch delete DSC template rule associations. Deleting this
  resource will not restore the deleted rule associations or undo the delete action, but will only remove the
  resource information from the tf state file.

## Example Usage

```hcl
variable "template_id" {
  type = string
}

variable "rule_ids" {
  type = string
}

resource "huaweicloud_dsc_batch_delete_template_rule" "test" {
  template_id = var.template_id
  rule_ids    = var.rule_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_id` - (Required, String, NonUpdatable) Specifies the scan template ID.

* `rule_ids` - (Required, String, NonUpdatable) Specifies the rule IDs to be deleted.
  If there are multiple IDs, separate them with commas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
