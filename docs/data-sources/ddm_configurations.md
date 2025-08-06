---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_configurations"
description: |-
  Use this data source to obtain parameter templates, including all default and custom parameter templates.
---

# huaweicloud_ddm_configurations

Use this data source to obtain parameter templates, including all default and custom parameter templates.

## Example Usage

```hcl
data "huaweicloud_ddm_configurations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - Indicates the list of DDM configurations.
  The [Configuration](#DdmConfigurations_Configuration) structure is documented below.

<a name="DdmConfigurations_Configuration"></a>
The `Configurations` block supports:

* `id` - Indicates the parameter template ID.

* `name` - Indicates the name of the parameter template.

* `description` - Indicates the description of the parameter template.

* `datastore_name` - Indicates the database type.

* `created` - Indicates the creation time, in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `updated` - Indicates the update time, in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `user_defined` - Indicates whether the parameter template is a custom template. Possible values: **false** (default
  template), **true** (custom template).
