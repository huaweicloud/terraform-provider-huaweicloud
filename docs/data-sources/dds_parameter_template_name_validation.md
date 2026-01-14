---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template_name_validation"
description: |-
  Use this data source to query whether the parameter template name exists.
---

# huaweicloud_dds_parameter_template_name_validation

Use this data source to query whether the parameter template name exists.

## Example Usage

```hcl
variable "param_template_name"  {}

data "huaweicloud_dds_parameter_template_name_validation" "test" {
  name = var.param_template_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the parameter template name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_existed` - Whether the parameter template name exists.
  The valid values are as follows:
  + **true**: Indicates the parameter template name exists.
  + **false**: Indicates the parameter template name does not exist.
