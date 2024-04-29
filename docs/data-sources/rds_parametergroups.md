---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroups"
description: ""
---

# huaweicloud_rds_parametergroups

Use this data source to get the list of RDS parametergroups.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_rds_parametergroups" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the parameter template name.

* `datastore_version_name` - (Optional, String) Specifies the database version name.

* `datastore_name` - (Optional, String) Specifies the database name.

* `user_defined` - (Optional, Bool) Specifies whether the parameter template is created by users.
  The options are as follows:
  + **false**: The parameter template is a default parameter template.
  + **true**: The parameter template is a custom template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - Lists the configurations.
The [configurations](#Rds_configurations) structure is documented below.

<a name="Rds_configurations"></a>
The `configurations` block supports:

* `id` - The parameter template ID.

* `name` - The parameter template name.

* `description` - The parameter template description.

* `datastore_version_name` - The database version name.

* `datastore_name` - The database name.

* `user_defined` - Whether the parameter template is created by users.
  The values can be **false** and **true**.

* `created_at` - The creation time of the configuration.

* `updated_at` - The latest update time of the configuration.
