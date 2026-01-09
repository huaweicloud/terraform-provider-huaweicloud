---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_distribution"
description: |-
  Manage an RDS distribution resource within HuaweiCloud.
---

# huaweicloud_rds_distribution

Manage an RDS distribution resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_distribution" "instance" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds PostgreSQL SQL limit resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS PostgreSQL instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID which is same as `instance_id`.

* `status` - Indicates the status.

* `distributor_instance_name` - Indicates the distributor instance name.

## Import

The RDS distribution can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_distribution.test <id>
```
