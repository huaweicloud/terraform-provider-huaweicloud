---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_backup_stop"
description: |-
  Manages the stop operation for an ongoing backup of a specified RDS instance within HuaweiCloud.
---

# huaweicloud_rds_backup_stop

Manages the stop operation for an ongoing backup of a specified RDS instance within HuaweiCloud.

-> This resource is a one-time action resource for operating the API.
Deleting this resource will not reverse the stop backup operation;
it will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_backup_stop" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance for which the
  ongoing backup is to be stopped.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.
