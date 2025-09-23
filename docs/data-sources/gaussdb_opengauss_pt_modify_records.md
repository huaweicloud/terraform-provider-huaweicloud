---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_pt_modify_records"
description: |-
  Use this data source to get the change history of a parameter template.
---

# huaweicloud_gaussdb_opengauss_pt_modify_records

Use this data source to get the change history of a parameter template.

## Example Usage

```hcl
variable "config_id"  {}

data "huaweicloud_gaussdb_opengauss_pt_modify_records" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `config_id` - (Required, String) Specifies the parameter template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - Indicates the parameter modify records.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `parameter_name` - Indicates the parameter name.

* `old_value` - Indicates the old parameter value.

* `new_value` - Indicates the new parameter value.

* `update_result` - Indicates the change status.
  The value can be **SUCCESS** or **FAILED**.

* `updated_at` - Indicates the modification time in the **yyyy-mm-ddThh:mm:ssZ** format.
