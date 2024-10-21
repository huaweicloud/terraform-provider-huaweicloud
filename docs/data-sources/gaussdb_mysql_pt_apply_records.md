---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_pt_apply_records"
description: |-
  Use this data source to get the application records of a parameter template.
---

# huaweicloud_gaussdb_mysql_pt_apply_records

Use this data source to get the application records of a parameter template.

## Example Usage

```hcl
variable "config_id"  {}

data "huaweicloud_gaussdb_mysql_pt_apply_records" "test" {
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

* `histories` - Indicates the parameter apply records.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `target_id` - Indicates the ID of the instance.

* `target_name` - Indicates the name of the instance.

* `apply_result` - Indicates the application result.

* `applied_at` - Indicates the application time.

* `error_code` - Indicates the error code.
