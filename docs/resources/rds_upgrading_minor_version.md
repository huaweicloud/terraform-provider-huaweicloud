---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance_minor_version_upgrade"
description: |-
  Manages an RDS instance minor version upgrade resource within HuaweiCloud.
---

# huaweicloud_rds_instance_minor_version_upgrade

Manages an RDS instance minor version upgrade resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_instance_minor_version_upgrade" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `is_delayed` - (Optional, Bool, NonUpdatable) Specifies whether the upgrade is delayed to the maintenance window.
  + **true**: Specifies the upgrade is delayed and performed within the maintenance window.
  + **false**: Specifies the upgrade is performed immediately. This is the default value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
