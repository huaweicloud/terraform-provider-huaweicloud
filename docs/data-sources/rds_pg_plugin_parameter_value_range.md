---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_plugin_parameter_value_range"
description: |-
  Use this data source to get the list of RDS PostgreSQL plugin parameter value range.
---

# huaweicloud_rds_pg_plugin_parameter_value_range

Use this data source to get the list of RDS PostgreSQL plugin parameter value range.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_pg_plugin_parameter_value_range" "test" {
  instance_id = var.instance_id
  name        = "shared_preload_libraries"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `name` - (Required, String) Specifies the parameter name. Currently only **shared_preload_libraries** is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `restart_required` - Indicates whether a reboot is required.

* `default_values` - Indicates the list of default parameter values. The default plugin will be loaded when the instance
  is created and can not be unloaded.

* `values` - Indicates the list of parameter values.
