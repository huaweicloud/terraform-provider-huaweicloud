---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_batch_update_template_rules_classification"
description: |-
  Use this resource to batch update template rules classification mapping within HuaweiCloud.
---

# huaweicloud_dsc_batch_update_template_rules_classification

Use this resource to batch update template rules classification mapping within HuaweiCloud.

-> This resource is a one-time action resource used to batch update DSC template rules classification mapping.
  Deleting this resource will not restore the updated classification mapping or undo the update action, but will only
  remove the resource information from the tf state file.

## Example Usage

```hcl
variable "template_id" {
  type = string
}

variable "classification_id" {
  type = string
}

variable "rule_id_list" {
  type = list(string)
}

resource "huaweicloud_dsc_batch_update_template_rules_classification" "test" {
  template_id       = var.template_id
  classification_id = var.classification_id
  rule_id_list      = var.rule_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `template_id` - (Required, String, NonUpdatable) Specifies the scan template ID.

* `classification_id` - (Optional, String, NonUpdatable) Specifies the classification ID for batch updating
  the rule classification.

* `rule_id_list` - (Optional, List, NonUpdatable) Specifies the list of rule IDs to be batch updated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
