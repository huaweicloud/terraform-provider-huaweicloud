---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_instances"
description: |-
  Use this data source to get the list of exclusive instances within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_instances

Use this data source to get the list of exclusive instances within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "exclusive_cluster_name" {}

data "huaweicloud_dataarts_dataservice_instances" "test" {
  workspace_id = var.workspace_id
  name         = var.exclusive_cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the exclusive clusters are located.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies ID of the workspace to which the exclusive clusters belong.

* `name` - (Optional, String) Specifies the exclusive cluster name to be queried.

* `create_user` - (Optional, String) Specifies the creator of the exclusive cluster to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `instances` - All exclusive clusters that match the filter parameters.  
  The [instances](#dataservice_instances_elem) structure is documented below.

<a name="dataservice_instances_elem"></a>
The `instances` block supports:

* `id` - The ID of the exclusive cluster, in UUID format.

* `name` - The name of the exclusive cluster.

* `description` - The description of the exclusive cluster.

* `external_address` - The external IP address of the exclusive cluster.

* `intranet_address` - The intranet IP address of the exclusive cluster.

* `intranet_address_ipv6` - The intranet IPv6 address of the exclusive cluster.

* `public_zone_id` - The public zone ID of the exclusive cluster.

* `public_zone_name` - The public zone name of the exclusive cluster.

* `private_zone_id` - The private zone ID of the exclusive cluster.

* `private_zone_name` - The private zone name of the exclusive cluster.

* `enterprise_project_id` - The enterprise project ID to which the exclusive cluster belongs.

* `created_at` - The create time of the exclusive cluster, in RFC3339 format.

* `create_user` - The create user of the exclusive cluster.

* `current_namespace_publish_api_num` - The number of the published APIs in the current namespace.

* `all_namespace_publish_api_num` - The number of the published APIs.

* `api_publishable_num` - The API quota of the exclusive cluster.

* `deletable` - Whether the exclusive cluster can be deleted.

* `status` - The status of the exclusive cluster.

* `flavor` - The flavor of the exclusive cluster.  
  The [flavor](#dataservice_instance_flavor_attr) structure is documented below.

* `gateway_version` - The version of the exclusive cluster.

* `availability_zone` - The availability zone where the exclusive cluster is located.

* `vpc_id` - The VPC ID to which the exclusive cluster belongs.

* `subnet_id` - The subnet ID to which the exclusive cluster belongs.

* `security_group_id` - The security group ID associated to the exclusive cluster.

* `node_num` - The node number of the exclusive cluster.

* `nodes` - The list of instance nodes.
  The [nodes](#dataservice_instance_nodes_attr) structure is documented below.

<a name="dataservice_instance_flavor_attr"></a>
The `flavor` block supports:

* `id` - The flavor ID.

* `name` - The flavor name.

* `disk_size` - The number of the disk size.

* `vcpus` - The number of CPU cores in the flavor.

* `memory` - The memory size in the flavor, in GB.

<a name="dataservice_instance_nodes_attr"></a>
The `nodes` block supports:

* `id` - The node ID.

* `name` - The node name.

* `private_ip` - The private IP address of the node.

* `status` - The status of the node.

* `created_at` - The create time of the node, in RFC3339 format.

* `create_user` - The create user of the node.

* `gateway_version` - The version of the node.
