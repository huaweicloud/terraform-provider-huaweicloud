---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_read_strategy"
description: ""
---

# huaweicloud_ddm_instance_read_strategy

Use this resource to set a DDM instance read strategy within HuaweiCloud.

-> **NOTE:** Deleting read strategy is not supported. If you destroy a resource of read strategy,
the read strategy is only removed from the state, but it remains in the cloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "rds_instance_id" {}
variable "rds_read_replica_instance_id" {}

resource "huaweicloud_ddm_instance_read_strategy" "test" {
  instance_id = var.instance_id

  read_weights {
    db_id  = var.rds_instance_id
    weight = 70
  }

  read_weights {
    db_id  = var.rds_read_replica_instance_id
    weight = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DDM instance.
  Changing this creates a new resource.

* `read_weights` - (Required, List) Specifies the list of read weights of the primary DB instance
  and its read replicas. The [read_weights](#ddm_read_weights) structure is documented below.

<a name="ddm_read_weights"></a>
The `read_weights` block supports:

* `db_id` - (Required, String) Specifies the ID of the DB instance associated with the DDM schema.

* `weight` - (Required, Int) Specifies read weight of the DB instance associated with the DDM schema.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the DDM instance ID.
