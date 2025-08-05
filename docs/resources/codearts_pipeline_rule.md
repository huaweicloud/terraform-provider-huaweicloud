---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_rule"
description: |-
  Manages a CodeArts pipeline rule resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_rule

Manages a CodeArts pipeline rule resource within HuaweiCloud.

## Example Usage

### create rule with metrics query from data source huaweicloud_codearts_pipeline_plugin_metrics

```hcl
variable "name" {}
variable "type" {}
variable "plugin_name" {}
variable "plugin_version" {}

data "huaweicloud_codearts_pipeline_plugin_metrics" "test" {
  plugin_name = var.plugin_name
  version     = var.plugin_version
}

locals {
  output_value = jsondecode(data.huaweicloud_codearts_pipeline_plugin_metrics.test.metrics[0].output_value)
}

resource "huaweicloud_codearts_pipeline_rule" "test" {
  name           = var.name
  type           = var.type
  layout_content = "layout_content"
  plugin_id      = var.plugin_name
  plugin_name    = var.plugin_name
  plugin_version = var.plugin_version

  dynamic "content" {
    for_each = local.output_value

    content {
      can_modify_when_inherit = true
      group_name              = content.value.group_name

      dynamic "properties" {
        for_each = content.value.properties

        content {
          is_valid   = true
          name       = properties.value.desc
          key        = properties.value.key
          value      = properties.value.value
          operator   = "="
          type       = "judge"
          value_type = "float"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the rule name.

* `type` - (Required, String, NonUpdatable) Specifies the rule type.
  Valid values are **Build**, **Gate**, **Deploy**, **Test** and **Normal**.

* `layout_content` - (Required, String, NonUpdatable) Specifies the layout content.

* `content` - (Required, List) Specifies the rule attribute group list.
  The [content](#block--content) structure is documented below.

* `plugin_id` - (Optional, String) Specifies the plugin ID.

* `plugin_name` - (Optional, String) Specifies the plugin name.

* `plugin_version` - (Optional, String) Specifies the plugin version.

<a name="block--content"></a>
The `content` block supports:

* `group_name` - (Required, String) Specifies the group name.

* `properties` - (Required, List) Specifies the rule attribute list.
  The [properties](#block--content--properties) structure is documented below.

* `can_modify_when_inherit` - (Optional, Bool) Specifies whether thresholds of an inherited policy can be modified.
  Default to **false**.

<a name="block--content--properties"></a>
The `properties` block supports:

* `key` - (Required, String) Specifies the attribute key.

* `name` - (Required, String) Specifies the display name.

* `type` - (Required, String) Specifies the type.

* `value` - (Required, String) Specifies the attribute value.

* `value_type` - (Required, String) Specifies the value type.

* `is_valid` - (Optional, Bool) Specifies wether the property is valid. Default to **false**.

* `operator` - (Optional, String) Specifies the comparison operators.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `version` - Indicates the rule version.

* `pipeline_count` - Indicates the number of pipelines.

* `project_count` - Indicates the number of projects.

* `rule_set_count` - Indicates the number of policies.

* `creator` - Indicates the creator.

* `create_time` - Indicates the create time.

* `updater` - Indicates the updater.

* `update_time` - Indicates the update time.

## Import

The rule can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_rule.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `layout_content`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the rule, or the resource definition should be updated to
align with the rule. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_rule" "test" {
    ...

  lifecycle {
    ignore_changes = [
      layout_content,
    ]
  }
}
```
