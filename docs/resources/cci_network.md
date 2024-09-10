---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cci_network"
description: ""
---

# huaweicloud_cci_network

Manages a CCI Network resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace_name" {}
variable "network_name" {}
variable "vpc_network_id" {}
variable "security_group_id" {}

resource "huaweicloud_cci_network" "test" {
  namespace         = var.namespace_name
  name              = var.network_name
  network_id        = var.vpc_network_id
  security_group_id = var.security_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCI network.
  If omitted, the provider-level region will be used. Changing this will create a new CCI network resource.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone (AZ) to which the CCI network
  belongs. Changing this will create a new CCI network resource.

* `namespace` - (Required, String, ForceNew) Specifies the namespace to logically divide your cloud container instances
  into different group. Changing this will create a new CCI network resource.

* `name` - (Required, String, ForceNew) Specifies an unique name of the CCI network resource.
  The name can contain a maximum of `200` characters, which may consist of lowercase letters, digits and hyphens (-).
  The name must start and end with a lowercase letter or digit. Changing this will create a new CCI network resource.

* `security_group_id` - (Required, String, ForceNew) Specifies a security group ID to which the CCI network belongs to.
  Changing this will create a new CCI network resource.

* `network_id` - (Required, String, ForceNew) Specifies a network ID of the VPC subnet which the CCI network belongs to.
  Changing this will create a new CCI network resource.

  ->**NOTE:** Namespace selected enterprise projects are different from Subnet (VPC) owned enterprise projects, and the
  namespaces created may not work correctly for permission reasons.
  And if too few IP addresses are available, the workloads may fail to function properly.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID, which is network name.

* `vpc_id` - VPC ID which the subnet and CCI network belongs to.

* `subnet_id` - IPv4 subnet ID.

* `cidr` - The network segment on which the subnet resides.

* `status` - The CCI network status, including **Initializing**, **Pending** and **Active**.

## Import

Networks can be imported using their `namespace` and `id`, separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_cci_network.test <namespace>/<id>
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.
