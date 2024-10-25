---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_cluster"
description: ""
---

# huaweicloud_cdm_cluster

Manages CDM cluster resource within HuaweiCloud.

## Example Usage

### create a cdm cluster

```hcl
variable "name" {}
variable "flavor_id" {}
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_cdm_cluster" "cluster" {
  name              = var.name
  availability_zone = var.availability_zone
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  subnet_id         = var.subnet_id
  vpc_id            = var.vpc_id
  security_group_id = var.secgroup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the cluster resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies cluster name. Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies available zone.
  Changing this parameter will create a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies flavor id. Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies subnet ID. Changing this parameter will create a new resource.

* `security_group_id` - (Required, String, ForceNew) Specifies security group ID.
 Changing this parameter will create a new resource.

* `version` - (Optional, String, ForceNew) Specifies the cluster version. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id.

* `is_auto_off` - (Optional, Bool, ForceNew) Specifies Whether to enable auto shutdown. The auto shutdown and scheduled
 startup/shutdown functions cannot be enabled at the same time. When auto shutdown is enabled, if no job is running in
  the cluster and no scheduled job is created, a cluster will be automatically shut down 15 minutes after it starts
   running to reduce costs. The default value is `false`. Changing this parameter will create a new resource.

* `schedule_boot_time` - (Optional, String, ForceNew) Specifies time for scheduled startup of a CDM cluster.
 The CDM cluster starts at this time every day. The scheduled startup/shutdown and auto shutdown function cannot be
  enabled at the same time. The time format is `hh:mm:ss`. Changing this parameter will create a new resource.

* `schedule_off_time` - (Optional, String, ForceNew) Specifies time for scheduled shutdown of a CDM cluster.
 The system shuts down directly at this time every day without waiting for unfinished jobs to complete.
 The scheduled startup/shutdown and auto shutdown function cannot be enabled at the same time.
  The time format is `hh:mm:ss`. Changing this parameter will create a new resource.

* `email` - (Optional, List) Specifies email address for receiving notifications when a table/file migration
 job fails or an EIP exception occurs. The max number is 20.

* `phone_num` - (Optional, List) Specifies phone number for receiving notifications when a table/file
 migration job fails or an EIP exception occurs. The max number is 20.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created` - Create time. The format is: `YYYY-MM-DDThh:mm:ss`.

* `status` - Status.

* `publid_ip` - Public ip.

* `public_endpoint` - EIP bound to the cluster.

* `flavor_name` - The flavor name. Format is `cdm.<flavor_type>`

* `instances` - Instance list. Structure is documented below.

The `instances` block contains:

* `id` - Instance ID.

* `name` - Instance name.

* `private_ip` - Private IP.

* `public_ip` - Public IP.

* `manage_ip` - Management IP address.

* `traffic_ip` - Traffic IP.

* `role` - Instance role.

* `type` - Instance type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `delete` - Default is 10 minutes.

## Import

Clusters can be imported by `id`. For example,

```bash
terraform import huaweicloud_cdm_cluster.test b11b407c-e604-4e8d-8bc4-92398320b847
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `email` and `phone_num`.
 It is generally recommended running `terraform plan` after importing a cluster.
 You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cdm_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      email, phone_num,
    ]
  }
}
```
