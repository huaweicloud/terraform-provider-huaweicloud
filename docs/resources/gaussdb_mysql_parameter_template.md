---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_parameter_template"
description: ""
---

# huaweicloud_gaussdb_mysql_parameter_template

Manages a GaussDB MySQL parameter template resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
   name = "test_mysql_parameter_template"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the parameter template name. The template name can contain 1 to 64 characters.
  Only letters (case-sensitive), digits, hyphens (-), underscores (_), and periods (.) are allowed.

* `description` - (Optional, String) Specifies the parameter template description. The description can consist of
  up to 256 characters, and cannot contain the carriage return characters or special characters (!<"='>&).

* `datastore_engine` - (Optional, String, ForceNew) Specifies the DB engine. Currently, only **gaussdb-mysql** is supported.

  Changing this parameter will create a new resource.

* `datastore_version` - (Optional, String, ForceNew) Specifies the DB version.

  Changing this parameter will create a new resource.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can specify parameter values based on a default parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time in the "yyyy-MM-ddTHH:mm:ssZ" format.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

* `updated_at` - Indicates the update time in the "yyyy-MM-ddTHH:mm:ssZ" format.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

## Import

The GaussDB Mysql parameter template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_parameter_template.test <id>
```
