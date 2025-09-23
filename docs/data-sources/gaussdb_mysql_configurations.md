---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_configurations"
description: |-
  Use this data source to get the list of parameter templates.
---

# huaweicloud_gaussdb_mysql_configurations

Use this data source to get the list of parameter templates.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_configurations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - Indicates the list of parameter templates.

  The [configurations](#configurations_struct) structure is documented below.

<a name="configurations_struct"></a>
The `configurations` block supports:

* `id` - Indicates the ID of the parameter template.

* `name` - Indicates the  name of the parameter template.

* `user_defined` - Indicates whether the parameter template is a custom template.
  + **false**: default parameter template.
  + **true**: custom template.

* `description` - Indicates the description of parameter template.

* `datastore_name` - Indicates the engine name.

* `datastore_version_name` - Indicates the engine version.

* `created_at` - Indicates the creation time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `updated_at` - Indicates the update time in the **yyyy-MM-ddTHH:mm:ssZ** format.
