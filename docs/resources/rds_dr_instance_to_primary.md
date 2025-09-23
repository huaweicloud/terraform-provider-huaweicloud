---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_dr_instance_to_primary"
description: |-
  Manages an RDS DR instance to primary DB instance resource within HuaweiCloud.
---

# huaweicloud_rds_dr_instance_to_primary

Manages an RDS DR instance to primary DB instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_dr_instance_to_primary" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS DR instance to primary resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS DR instance.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the DR instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
