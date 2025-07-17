---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_modules"
description: |-
  Use this data source to a list of CodeArts pipeline modules.
---

# huaweicloud_codearts_pipeline_modules

Use this data source to a list of CodeArts pipeline modules.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_modules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the module name.

* `tags` - (Optional, List) Specifies the tags.

* `product_line` - (Optional, String) Specifies the product line.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `modules` - Indicates the module list.
  The [modules](#attrblock--modules) structure is documented below.

<a name="attrblock--modules"></a>
The `modules` block supports:

* `id` - Indicates the module ID.

* `name` - Indicates the module name.

* `module_id` - Indicates the module ID.

* `base_url` - Indicates the module base URL.

* `description` - Indicates the module description.

* `location` - Indicates the endpoint.

* `manifest_version` - Indicates the summary version.

* `properties` - Indicates the properties.

* `properties_list` - Indicates the properties list.

* `publisher` - Indicates the publisher.

* `tags` - Indicates the tags.

* `type` - Indicates the module type.

* `url_relative` - Indicates the extension URL.

* `version` - Indicates the module version.
