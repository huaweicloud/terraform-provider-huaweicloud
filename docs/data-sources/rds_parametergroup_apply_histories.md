---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_apply_histories"
description: |-
  Use this data source to get the query application records of a parameter template.
---

# huaweicloud_rds_parametergroup_apply_histories

Use this data source to get the query application records of a parameter template.

## Example Usage

```hcl
variable "config_id" {}
variable "instance_id" {}

resource "huaweicloud_rds_parametergroup_apply_histories" "test" {
  config_id   = var.config_id
  instance_id = var.instance_id
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

* `histories` - Indicates the list of historical application details.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `instance_id` - Indicates the ID of the instance to which the parameter template is applied.

* `instance_name` - Indicates the name of the instance to which the parameter template is applied.

* `apply_result` - Indicates the result of applying the parameter template.

* `apply_time` - Indicates the time when the parameter template is applied.

* `error_code` - Indicates the error code displayed when you submit the request to apply the parameter template.
