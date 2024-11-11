---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_pt_modification_records"
description: |-
  Use this data source to get the list of DDS parameter template modification records.
---

# huaweicloud_dds_pt_modification_records

Use this data source to get the list of DDS parameter template modification records.

## Example Usage

```hcl
variable "configuration_id"  {}

data "huaweicloud_dds_pt_modification_records" "test" {
  configuration_id = var.configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `configuration_id` - (Required, String) Specifies the ID of the parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - Indicates the modification records.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `parameter_name` - Indicates the parameter name.

* `new_value` - Indicates the new value.

* `old_value` - Indicates the old value.

* `updated_at` - Indicates the update time, in the **yyyy-mm-ddThh:mm:ssZ** format.
