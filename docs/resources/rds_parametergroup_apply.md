---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_apply"
description: |-
  Manages an RDS parameter group apply resource within HuaweiCloud.
---

# huaweicloud_rds_parametergroup_apply

Manages an RDS parameter group apply resource within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}
variable "instance_id" {}

resource "huaweicloud_rds_parametergroup_apply" "test" {
  config_id   = var.config_id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS binlog resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `config_id` - (Required, String, NonUpdatable) Specifies the parameter template ID.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID .
