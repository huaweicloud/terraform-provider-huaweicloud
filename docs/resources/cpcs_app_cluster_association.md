---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_app_cluster_association"
description: |-
  Manages an association between a CPCS application and cluster within HuaweiCloud.
---

# huaweicloud_cpcs_app_cluster_association

Manages an association between a CPCS application and cluster within HuaweiCloud.

-> Currently, this resource is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "app_id" {}
variable "cluster_id" {}

resource "huaweicloud_cpcs_app_cluster_association" "test" {
  app_id     = var.app_id
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_id` - (Required, String, NonUpdatable) Specifies the ID of the application to be associated with the cluster.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to associate with the application.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (Randomly generated UUID).

* `cluster_name` - The name of the cluster.

* `app_name` - The name of the application.

* `vpc_name` - The name of the VPC where the application is located.

* `subnet_name` - The name of the subnet where the application is located.

* `cluster_server_type` - The type of the cluster server.

* `vpcep_address` - The address of the VPC endpoint.

* `update_time` - The update time of the association, UNIX timestamp in milliseconds.

* `create_time` - The creation time of the association, UNIX timestamp in milliseconds.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The CPCS application cluster association resource can be imported using the `app_id` and `cluster_id`,
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cpcs_app_cluster_association.test <app_id>/<cluster_id>
```
