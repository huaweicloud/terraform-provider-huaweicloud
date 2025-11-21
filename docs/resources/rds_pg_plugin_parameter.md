---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_plugin_parameter"
description: |-
  Manage an RDS PostgreSQL plugin parameter resource within HuaweiCloud.
---

# huaweicloud_rds_pg_plugin_parameter

Manage an RDS PostgreSQL plugin parameter resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_plugin_parameter" "test" {
  instance_id = var.instance_id
  name        = "shared_preload_libraries"
  values      = ["pg_stat_kcache"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds PostgreSQL plugin parameter resource.
  If omitted, the provider-level region will be used. Changing this will creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS instance.

* `name` - (Required, String, NonUpdatable) Specifies the name of the plugin parameter.

* `values` - (Required, List) Specifies the list of plugin parameter values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID. The format is `<instance_id>/<name>`.

* `restart_required` - Indicates whether a reboot is required.

* `default_values` - Indicates the default values of the plugin parameter.

## Import

The plugin parameter can be imported using the `instance_id` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_pg_plugin_parameter.test <instance_id>/<name>
```
