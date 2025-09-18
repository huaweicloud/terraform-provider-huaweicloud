---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_private_network_access_control"
description: |-
  Manages a SWR enterprise instance private network access control resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_private_network_access_control

Manages a SWR enterprise instance private network access control resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "description" {}

resource "huaweicloud_swr_enterprise_private_network_access_control" "test" {
  instance_id = var.instance_id
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID.

* `description` - (Optional, String, NonUpdatable) Specifies the description.

* `project_id` - (Optional, String, NonUpdatable) Specifies the project ID to which the VPC belongs.
  Default to the project ID of the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vpcep_endpoint_id` - Indicates the VPCEP endpoint ID.

* `endpoint_ip` - Indicates the endpoint IP.

* `status` - Indicates the endpoint status.

* `status_text` - Indicates the status text

* `project_name` - Indicates the project name to which the VPC belongs.

* `vpc_cidr` - Indicates the VPC cidr.

* `vpc_name` - Indicates the VPC name.

* `subnet_cidr` - Indicates the subnet cidr.

* `subnet_name` - Indicates the subnet name.

* `created_at` - Indicates the creation time of the private network access control.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The private network access control can be imported using `instance_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_private_network_access_control.test <instance_id>/<id>
```
