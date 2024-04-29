---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_associate"
description: ""
---

# huaweicloud_global_eip_associate

Associates a GEIP to a specified instance.

## Example Usage

```hcl
variable "global_eip_id" {}
variable "region" {}
variable "project_id" {}
variable "compute_instance_id" {}
variable "gc_bandwidth_name" {}

resource "huaweicloud_global_eip_associate" "test" {
  global_eip_id  = var.global_eip_id
  is_reserve_gcb = false

  associate_instance {
    region        = var.region
    project_id    = var.project_id
    instance_type = "ECS"
    instance_id   = var.compute_instance_id
  }

  gc_bandwidth {
    name        = var.gc_bandwidth_name
    charge_mode = "bwd"
    size        = 5
  }
}
```

## Argument Reference

The following arguments are supported:

* `global_eip_id` - (Required, String, ForceNew) Specifies the global EIP ID.
  Changing this creates a new resource.

* `is_reserve_gcb` - (Required, Bool) Specifies whether to reserve the GCB when the GEIP disassociates to the instance.

* `associate_instance` - (Required, List, ForceNew) Specifies the information of instance which the GEIP associates to.
  Changing this creates a new resource.
  The [associate_instance](#block--associate_instance) structure is documented below.

* `gc_bandwidth` - (Optional, List, ForceNew) Specifies the information of GCB which the GEIP associates to.
  Changing this creates a new resource.
  The [gc_bandwidth](#block--gc_bandwidth) structure is documented below.

<a name="block--associate_instance"></a>
The `associate_instance` block supports:

* `region` - (Required, String, ForceNew) Specifies the region of the instance.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID of the region.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `instance_type` - (Required, String, ForceNew) Specifies the instance type. Valid values are **ECS**, **PORT**,
  **NATGW** and **ELB**. If value is **ECS** or **PORT**, make sure the VPC associating with an internet gateway.
  Changing this creates a new resource.

* `service_id` - (Optional, String, ForceNew) Specifies the service ID.
  Changing this creates a new resource.

* `service_type` - (Optional, String, ForceNew) Specifies the service type.
  Changing this creates a new resource.

<a name="block--gc_bandwidth"></a>
The `gc_bandwidth` block supports:

* `id` - (Optional, String, ForceNew) Specifies the GCB ID which is existing.
  Changing this creates a new resource.

* `name` - (Optional, String, ForceNew) Specifies the GCB name. When `gc_bandwidth.id` is empty, it is **Required** for
  creating a new GCB. Changing this creates a new resource.

* `charge_mode` - (Optional, String, ForceNew) Specifies the GCB charge mode. When `gc_bandwidth.id` is empty, it is
  **Required** for creating a new GCB.

  Valid values are as follows:
  + **bwd**: Billed by bandwidth.
  + **95**: Billed by 95th percentile bandwidth.

  Changing this creates a new resource.

* `size` - (Optional, Int, ForceNew) Specifies the GCB size. When `gc_bandwidth.id` is empty, it is **Required** for
  creating a new GCB. If `gc_bandwidth.charge_mode` is **95**, the range is **100-300 Mbit/s**, otherwise, the range is
  **2-300 Mbit/s**. Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of GCB.
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of GCB.
  Changing this creates a new resource.

* `tags` - (Optional, Map, ForceNew) Specifies the tags of GCB.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. Same with the global EIP ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The global EIP association can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_global_eip_associate.test <id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `is_reserve_gcb`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_global_eip_associate" "test" {
    ...

  lifecycle {
    ignore_changes = [
      is_reserve_gcb,
    ]
  }
}
```
