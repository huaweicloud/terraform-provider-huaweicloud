---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_patch"
description: |-
  Use this data source to query the SQL patch information of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_patch

Use this data source to query the SQL patch information of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "sql_id" {}

data "huaweicloud_gaussdb_sql_patch" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  sql_id      = var.sql_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SQL patch.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String) Specifies the node ID of the GaussDB instance.

* `sql_id` - (Required, String) Specifies the SQL ID to query the patch information.

* `database_name` - (Optional, String) Specifies the database name. This parameter is optional in slow SQL scenarios and
  mandatory in other scenarios.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `patch_name` - The name of the SQL patch.

* `hint` - The SQL patch content (Hint text). Leave this field empty for ABORT SQL patches.

* `patch_status` - The status of the SQL patch.
  The valid values are as follows:
  + **enabled**: The SQL patch is effective.
  + **disabled**: The SQL patch is ineffective.
