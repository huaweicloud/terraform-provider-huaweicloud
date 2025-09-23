---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_primary_instance_dr_capability"
description: |-
  Manages RDS primary instance DR capability resource within HuaweiCloud.
---

# huaweicloud_rds_primary_instance_dr_capability

Manages RDS primary instance DR capability resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "target_instance_id" {}
variable "target_project_id" {}
variable "target_region" {}
variable "target_ip" {}
variable "target_subnet" {}

resource "huaweicloud_rds_primary_instance_dr_capability" "test" {
  instance_id        = var.instance_id
  target_instance_id = var.target_instance_id
  target_project_id  = var.target_project_id
  target_region      = var.target_region
  target_ip          = var.target_ip
  target_subnet      = var.target_subnet
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `target_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DR instance.

* `target_project_id` - (Required, String, NonUpdatable) Specifies the project ID of the tenant to which the DR instance
  belongs.

* `target_region` - (Required, String, NonUpdatable) Specifies the ID of the region where the DR instance resides.

* `target_ip` - (Required, String, NonUpdatable) Specifies the data virtual IP address (VIP) of the DR instance.

* `target_subnet` - (Required, String, NonUpdatable) Specifies the subnet IP address of the DR instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the DR configuration status.

* `time` - Indicates the DR configuration time.

* `build_process` - Indicates the process for configuring disaster recovery (DR). The value can be:
  + **master**: process of configuring DR capability for the primary instance
  + **slave**: process of configuring DR for the DR instance

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The RDS primary instance dr capability can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_primary_instance_dr_capability.test <id>
```
