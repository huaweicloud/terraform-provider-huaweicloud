---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_applied_instances"
description: |-
  Use this data source to get the DB instances that a parameter template is applied to.
---

# huaweicloud_rds_parametergroup_applied_instances

Use this data source to get the DB instances that a parameter template is applied to.

## Example Usage

```hcl
variable "config_id" {}

resource "huaweicloud_rds_parametergroup_applied_instances" "test" {
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

* `entities` - Indicates the list of instances.

  The [entities](#entities_struct) structure is documented below.

* `instance_count_limit` - Indicates the limit on the number of instances.

<a name="entities_struct"></a>
The `entities` block supports:

* `entity_id` - Indicates the instance ID.

* `entity_name` - Indicates the instance name.
