---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_templates"
description: |-
  Use this data source to get the list of SecMaster component templates.
---

# huaweicloud_secmaster_component_templates

Use this data source to get the list of SecMaster component templates.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}

data "huaweicloud_secmaster_component_templates" "test" {
  workspace_id = var.workspace_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `component_id` - (Required, String) Specifies the component ID.

* `file_type` - (Optional, String) Specifies the file type. The valid values are **JVM**, **LOG4J2**, and **YML**.

* `sort_key` - (Optional, String) Specifies the attribute fields for sorting.

* `sort_dir` - (Optional, String) Specifies the sorting order. Supported values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The component templates list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `version` - The version.

* `component_id` - The component ID.

* `component_name` - The component name.

* `param` - The param.

* `file_type` - The configuration file type. Can be used to identify the type of configuration file.
  The valid values are **JVM**, **LOG4J2**, and **YML**.

* `file_name` - The file name.

* `file_path` - The file path.
