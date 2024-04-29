---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_code_table_values"
description: ""
---

# huaweicloud_dataarts_architecture_code_table_values

Manages a DataArts Architecture code table values resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "table_id" {}
variable "field_id" {}
variable "field_ordinal" {}
variable "field_name" {}
variable "field_code" {}
variable "field_type" {}

resource "huaweicloud_dataarts_architecture_code_table_values" "test" {
  workspace_id  = var.workspace_id
  table_id      = var.table_id
  field_id      = var.field_id
  field_ordinal = var.field_ordinal
  field_name    = var.field_name
  field_code    = var.field_code
  field_type    = var.field_type

  values {
    value = "1"  
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of DataArts Studio workspace.
  Changing this parameter will create a new resource.

* `table_id` - (Required, String, ForceNew) Specifies the ID of the code table.
  Changing this parameter will create a new resource.

* `field_id` - (Required, String, ForceNew) Specifies the ID of the code table field.
  Changing this parameter will create a new resource.

* `field_ordinal` - (Required, Int, ForceNew) Specifies the ordinal of the code table field.
  Changing this parameter will create a new resource.

* `field_name` - (Required, String, ForceNew) Specifies the name of the code table field.
  Changing this parameter will create a new resource.

* `field_code` - (Required, String, ForceNew) Specifies the code of the code table field.
  Changing this parameter will create a new resource.

* `field_type` - (Required, String, ForceNew) Specifies the type of the code table field.
  Changing this parameter will create a new resource.

* `values` - (Required, List) Specifies the values list of the code table field.
  The [values](#Values) structure is documented below.

<a name="Values"></a>
The `values` block supports:

* `value` - (Required, String) Specifies the value of a field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `values` - The values list of the code table field.
  The [values](#Values_Attr) structure is documented below.

<a name="Values_Attr"></a>
The `values` block supports:

* `id` - The ID of a value.

* `ordinal` - The ordinal of a value.
