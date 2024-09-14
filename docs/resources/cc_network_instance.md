---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_network_instance"
description: ""
---

# huaweicloud_cc_network_instance

Manages a network instance resource within HuaweiCloud.  

Load VPCs or virtual gateways to the cloud connection. If you load virtual gateways, your on-premises data center
can access multiple VPCs to build a hybrid cloud.  

Each network instance can be loaded onto only one cloud connection.

## Example Usage

```hcl
variable "cloud_connection_id" {}
variable "vpc_instance_id" {}
variable "vpc_project_id" {}
variable "vpc_region_id" {}
variable "cidr" {}

resource "huaweicloud_cc_network_instance" "test" {
  type                = "vpc"
  cloud_connection_id = var.cloud_connection_id
  instance_id         = var.vpc_instance_id
  project_id          = var.vpc_project_id
  region_id           = var.vpc_region_id
  cidrs = [
    var.cidr
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Type of the network instance to be loaded to the cloud connection.  
  The options are as follows:
    + **vpc**: Virtual Private Cloud.
    + **vgw**: virtual gateway.

  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) ID of the VPC or virtual gateway to be loaded to the cloud connection.

  Changing this parameter will create a new resource.

* `cidrs` - (Required, List) List of routes advertised by the VPC or virtual gateway.

* `project_id` - (Required, String, ForceNew) Project ID of the VPC or virtual gateway.

  Changing this parameter will create a new resource.

* `region_id` - (Required, String, ForceNew) Region ID of the VPC or virtual gateway.

  Changing this parameter will create a new resource.

* `cloud_connection_id` - (Required, String, ForceNew) Cloud connection ID.

  Changing this parameter will create a new resource.

* `name` - (Optional, String) The network instance name.  
  The name can contain `1` to `64` characters, only English letters, Chinese characters, digits, hyphens (-),
  underscores (_) and dots (.).

* `description` - (Optional, String) The description about the network instance.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `instance_domain_id` - (Optional, String, ForceNew) Account ID of the VPC or virtual gateway.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_id` - Account ID.

* `status` - Network instance status.  
  The options are as follows:
    + **ACTIVE**: The network instance is available.

## Import

The network instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_network_instance.test 0ce123456a00f2591fabc00385ff1234
```
