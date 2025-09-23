---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_copy"
description: |-
  Manages an RDS parameter group copy resource within HuaweiCloud.
---

# huaweicloud_rds_parametergroup_copy

Manages an RDS parameter group copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}

resource "huaweicloud_rds_parametergroup_copy" "test" {
  config_id = var.config_id
  name      = "test_name"
  
  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `config_id` - (Required, String, NonUpdatable) Specifies the source parameter group ID.

* `name` - (Required, String) Specifies the parameter group name. It contains a maximum of 64 characters.

* `description` - (Optional, String) Specifies the parameter group description. It contains a maximum of 256 characters
  and cannot contain the following special characters:>!<"&'= the value is left blank by default.

* `values` - (Optional, Map) Specifies the parameter group values key/value pairs defined by users based on the default
  parameter groups.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `datastore` - Indicates the database object.
  The [datastore](#datastore_struct) structure is documented below.

* `configuration_parameters` - Indicates the parameter configuration defined by users based on the default parameters groups.
  The [configuration_parameters](#configuration_parameters_struct) structure is documented below.

* `created_at` - Indicates the creation time, in UTC format.

* `updated_at` - Indicates the last update time, in UTC format.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.

<a name="configuration_parameters_struct"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `restart_required` - Indicates whether a restart is required.

* `readonly` - Indicates whether the parameter is read-only.

* `value_range` - Indicates the parameter value range.

* `type` - Indicates the parameter type.

* `description` - Indicates the parameter description.

## Import

The RDS parameter group copy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_parametergroup_copy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `config_id` and`values`. It is generally
recommended running `terraform plan` after importing the RDS parameter group copy. You can then decide if changes should
be applied to the RDS parameter group copy, or the resource definition should be updated to align with the RDS parameter
group copy. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rds_parametergroup_copy" "test" {
    ...

  lifecycle {
    ignore_changes = [
      config_id, values,
    ]
  }
}
```
