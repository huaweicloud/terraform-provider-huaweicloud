---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_parameter_templates"
description: |-
  Use this data source to get the list of GaussDB OpenGauss parameter templates.
---

# huaweicloud_gaussdb_opengauss_parameter_templates

Use this data source to get the list of GaussDB OpenGauss parameter templates.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_parameter_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - Indicates the parameter template information

  The [configurations](#configurations_struct) structure is documented below.

<a name="configurations_struct"></a>
The `configurations` block supports:

* `id` - Indicates the parameter template ID, which is the unique ID of a parameter template.

* `name` - Indicates the parameter template name.

* `description` - Indicates the parameter template description.

* `datastore_version` - Indicates the engine version.

* `datastore_name` - Indicates the engine name.

* `ha_mode` - Indicates the instance type.

* `user_defined` - Indicates whether the parameter template is a custom template.
  The value can be:
  + **false**: The parameter template is a default template.
  + **true**: The parameter template is a custom template.

* `created_at` - Indicates the creation time in the **yyyy-MM-dd HH:mm:ss** format.

* `updated_at` - Indicates the update time in the **yyyy-MM-dd HH:mm:ss** format.
