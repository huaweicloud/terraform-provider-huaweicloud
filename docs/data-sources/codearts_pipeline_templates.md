---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_templates"
description: |-
  Use this data source to get a list of CodeArts pipeline templates.
---

# huaweicloud_codearts_pipeline_templates

Use this data source to get a list of CodeArts pipeline templates.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the template name.

* `is_system` - (Optional, Bool) Specifies whether the template is a system template. Default to **false**.

* `language` - (Optional, String) Specifies the template language. Value can be **Java**, **Python**, **Node.js**,
  **Go**, **.NET**, **CPP**, **PHP**, **other**, and **none**.

* `sort_dir` - (Optional, String) Specifies the sorting sequence. Value can be **asc** and **desc**.

* `sort_key` - (Optional, String) Specifies the sorting attribute. Value can be **name** and **create_time**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - Indicates the template list.
  The [templates](#attrblock--templates) structure is documented below.

<a name="attrblock--templates"></a>
The `templates` block supports:

* `id` - Indicates the template ID.

* `name` - Indicates the template name.

* `description` - Indicates the template description.

* `is_favorite` - Indicates whether it is a favorite template.

* `is_show_source` - Indicates whether to display the pipeline source.

* `is_system` - Indicates whether the template is a system template.

* `language` - Indicates the template language.

* `manifest_version` - Indicates the manifest version.

* `icon` - Indicates the template icon.

* `stages` - Indicates the stage running information.
  The [stages](#attrblock--templates--stages) structure is documented below.

* `create_time` - Indicates the creation time.

* `creator_id` - Indicates the creator.

* `creator_name` - Indicates the creator name.

* `update_time` - Indicates the last update time.

* `updater_id` - Indicates the last updater.

<a name="attrblock--templates--stages"></a>
The `stages` block supports:

* `name` - Indicates the stage name.

* `sequence` - Indicates the serial number.
