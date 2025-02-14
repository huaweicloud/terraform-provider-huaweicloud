---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_sql_templates"
description: |-
  Use this data source to get the SQL templates of a specified node.
---

# huaweicloud_gaussdb_opengauss_sql_templates

Use this data source to get the SQL templates of a specified node.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_gaussdb_opengauss_sql_templates" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of a GaussDB OpenGauss instance.

* `node_id` - (Required, String) Specifies the ID of a GaussDB OpenGauss instance node.

* `sql_model` - (Optional, String) Specifies the SQL template.
  The value can contain only uppercase letters, lowercase letters, underscores (_), digits, spaces,
  and the following special characters $*?=+;()><,.".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_limit_sql_model_list` - Indicates the information about the SQL template for SQL throttling.

  The [node_limit_sql_model_list](#node_limit_sql_model_list_struct) structure is documented below.

<a name="node_limit_sql_model_list_struct"></a>
The `node_limit_sql_model_list` block supports:

* `sql_id` - Indicates the SQL ID of the throttling task.

* `sql_model` - Indicates the SQL template of the throttling task.
