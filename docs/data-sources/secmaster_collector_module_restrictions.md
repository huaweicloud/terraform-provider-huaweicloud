---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_module_restrictions"
description: |-
  Use this data source to get the list of collector module restrictions.
---

# huaweicloud_secmaster_collector_module_restrictions

Use this data source to get the list of collector module restrictions.

## Example Usage

```hcl
variable "workspace_id" {}
variable "template_ids" {
  type = list(string)
}

data "huaweicloud_secmaster_collector_module_restrictions" "test" {
  workspace_id = var.workspace_id
  template_ids = var.template_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `template_ids` - (Required, List) Specifies the template IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `module_restrictions` - The module restrictions information.

  The [module_restrictions](#module_restrictions_struct) structure is documented below.

<a name="module_restrictions_struct"></a>
The `module_restrictions` block supports:

* `template_id` - The template ID.

* `fields` - The fields information.

  The [fields](#template_fields_struct) structure is documented below.

<a name="template_fields_struct"></a>
The `fields` block supports:

* `default_value` - The default value.

* `description` - The description.

* `example` - The example.

* `name` - The rule name.

* `required` - Whether required.

* `restrictions` - The restrictions information.

  The [restrictions](#restrictions_struct) structure is documented below.

* `template_field_id` - The template field ID.

* `title` - The title.

* `type` - The rule type.

<a name="restrictions_struct"></a>
The `restrictions` block supports:

* `logic` - The logic condition.

* `title` - The title.

* `type` - The rule type.

* `value` - The rule name.
