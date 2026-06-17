---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_preprocess_rule"
description: |-
  Manages a SOC preprocess rule resource within HuaweiCloud.
---

# huaweicloud_secmaster_soc_preprocess_rule

Manages a SOC preprocess rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "mapping_id" {}
variable "mapper_id" {}
variable "mapper_type_id" {}

resource "huaweicloud_secmaster_soc_preprocess_rule" "test" {
  workspace_id = var.workspace_id
  mapping_id   = var.mapping_id

  preprocess_rules {
    name           = "test_rule"
    mapper_id      = var.mapper_id
    mapper_type_id = var.mapper_type_id
    action         = "drop"
    expression     = "expression_content"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the preprocess rule
  belongs.

* `mapping_id` - (Required, String, NonUpdatable) Specifies the mapping ID of the preprocess rule.

* `preprocess_rules` - (Required, List) Specifies the list of preprocess rules.
  The [preprocess_rules](#soc_preprocess_rule_preprocess_rules) structure is documented below.

  -> Updating this field `preprocess_rules` will overwrite historical rules, and a new ID will be generated for each rule.

<a name="soc_preprocess_rule_preprocess_rules"></a>
The `preprocess_rules` block supports:

* `name` - (Optional, String) Specifies the name of the preprocess rule.

* `mapper_id` - (Optional, String) Specifies the mapper ID of the preprocess rule.

* `mapper_type_id` - (Optional, String) Specifies the mapper type ID of the preprocess rule.

* `action` - (Optional, String) Specifies the preprocess action. The value can be **drop** (discard).

* `expression` - (Optional, String) Specifies the expression of the preprocess rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the value of `mapping_id`.

* `data` - The list of preprocess rule details returned by the API.
  The [data](#soc_preprocess_rule_data) structure is documented below.

<a name="soc_preprocess_rule_data"></a>
The `data` block supports:

* `id` - The ID of the preprocess rule.

* `name` - The name of the preprocess rule.

* `project_id` - The project ID of the preprocess rule.

* `workspace_id` - The workspace ID of the preprocess rule.

* `mapping_id` - The mapping ID of the preprocess rule.

* `mapper_id` - The mapper ID of the preprocess rule.

* `mapper_type_id` - The mapper type ID of the preprocess rule.

* `action` - The preprocess action.

* `expression` - The expression of the preprocess rule.

* `creator_id` - The creator ID of the preprocess rule.

* `creator_name` - The creator name of the preprocess rule.

* `create_time` - The creation time of the preprocess rule.

* `update_time` - The latest update time of the preprocess rule.

* `modifier_id` - The modifier ID of the preprocess rule.

* `modifier_name` - The modifier name of the preprocess rule.

## Import

The SOC preprocess rule can be imported using the `workspace_id` and their `mapping_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_soc_preprocess_rule.test <workspace_id>/<mapping_id>
```
