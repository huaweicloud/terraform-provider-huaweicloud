---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_parameter_template"
description: |-
  Manages a GaussDB MySQL parameter template resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_parameter_template

Manages a GaussDB MySQL parameter template resource within HuaweiCloud.

## Example Usage

### create parameter template

```hcl
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" { 
  name = "test_mysql_parameter_template"
}
```

### replica parameter template from existed configuration

```hcl
variable "source_configuration_id" {}

resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name                    = "test_copy_from_configuration"
  source_configuration_id = var.source_configuration_id
}
```

### replica parameter template from existed instance

```hcl
variable "instance_id" {}
variable "instance_configuration_id" {}

resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name                      = "test_copy_from_instance"
  instance_id               = var.instance_id
  instance_configuration_id = var.instance_configuration_id
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

  -> **NOTE:** It is mandatory when `datastore_engine` is specified.

* `source_configuration_id` - (Optional, String, ForceNew) Specifies the source parameter template ID.

  Changing this parameter will create a new resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

  Changing this parameter will create a new resource.

* `instance_configuration_id` - (Optional, String, ForceNew) Specifies the parameter template ID of the GaussDB MySQL
  instance.

  Changing this parameter will create a new resource.

  -> **NOTE:** It is mandatory when `instance_id` is specified.

-> **NOTE:** 1. At most one of `datastore_engine`, `source_configuration_id` and `instance_id` can be specified.
  <br>2. If `source_configuration_id` is specified, then the resource will replicate the parameter template specified
  by `source_configuration_id`.
  <br>3. If `instance_id` is specified, then the resource will replicate the parameter template of the GaussDB MySQL
  instance specified by `instance_id`.
  <br>4. If `source_configuration_id` and `instance_id` are both not specified, then a new parameter template will be
  created directly.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can specify parameter values based on a default parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time in the **yyyy-MM-ddTHH:mm:ssZ** format.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

* `updated_at` - Indicates the update time in the **yyyy-MM-ddTHH:mm:ssZ** format.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

## Import

The GaussDB Mysql parameter template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_parameter_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `source_configuration_id`, `instance_id`
and `instance_configuration_id`. It is generally recommended running `terraform plan` after importing an GaussDB MySQL
parameter template. You can then decide if changes should be applied to the GaussDB MySQL parameter template, or the
GaussDB MySQL parameter template definition should be updated to align with the instance. Also, you can ignore changes
as below.

```hcl
resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_configuration_id, instance_id, instance_configuration_id,
    ]
  }
}
```
