---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_compare"
description: |-
  Manages an RDS parameter group compare resource within HuaweiCloud.
---

# huaweicloud_rds_parametergroup_compare

Manages an RDS parameter group compare resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_id" {}
variable "target_id" {}

resource "huaweicloud_rds_parametergroup_compare" "test" {
  source_id = var.source_id
  target_id = var.target_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `source_id` - (Required, String, NonUpdatable) Specifies the source parameter template ID.

* `target_id` - (Required, String, NonUpdatable) Specifies the target parameter template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `source_name` - Indicates the source parameter template name.

* `target_name` - Indicates the target parameter template name.

* `parameters` - Indicates the template parameter differences..
  The [parameters](#parameters_struct) structure is documented below.

<a name="parameters_struct"></a>
The `parameters` block supports:

* `name` - Indicates the parameter name.

* `source_value` - Indicates the parameter value in the source template.

* `target_value` - Indicates the parameter value in the target template.
