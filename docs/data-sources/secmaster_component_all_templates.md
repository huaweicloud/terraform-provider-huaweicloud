---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_all_templates"
description: |-
  Use this data source to get the list of component all templates.
---

# huaweicloud_secmaster_component_all_templates

Use this data source to get the list of component all templates.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_component_all_templates" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the component all templates belong.

* `search_value` - (Optional, String) Specifies the template name to search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of component all template details.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The ID of the component template.

* `component_id` - The component ID of the component template.

* `template_name` - The name of the component template.

* `task_config` - The task config of the component template.

* `create_time` - The creation time of the component template.

* `update_time` - The update time of the component template.

* `project_id` - The project ID of the component template.
