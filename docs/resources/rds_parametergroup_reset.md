---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup_reset"
description: |-
  Manages an RDS parameter group copy resource within HuaweiCloud.
---

# huaweicloud_rds_parametergroup_reset

Manages an RDS parameter group reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}

resource "huaweicloud_rds_parametergroup_reset" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `config_id` - (Required, String, NonUpdatable) Specifies the parameter template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `config_name` - Indicates the parameter template name.

* `need_restart` - Indicates whether a reboot is required.
