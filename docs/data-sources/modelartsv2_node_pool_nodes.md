---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_pool_nodes"
description: |-
  Use this data source to get node list under a specified node pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_pool_nodes

Use this data source to get node list under a specified node pool within HuaweiCloud.

## Example Usage

```hcl
variable "resource_pool_name" {}
variable "node_pool_name" {}

data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = var.resource_pool_name
  node_pool_name     = var.node_pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used.

* `resource_pool_name` - (Required, String) Specifies the resource pool name to which the node pool belongs.  

* `node_pool_name` - (Required, String) Specifies the node pool name to which the resource nodes belongs.

  -> If the node pool name does not specified when the resource pool creation, the ModelArts service will auto generate
     a name as this format: `{flavor}-default`, such as **modelarts.bm.npu.arm.8snt9b3.i6l.d.c004-default**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - All queried resource nodes under a specified node pool.  
  The [nodes](#v2modelarts_node_pool_nodes) structure is documented below.

<a name="v2modelarts_node_pool_nodes"></a>
The `nodes` block supports:

* `metadata` - The metadata information of the node.  
  The [metadata](#v2modelarts_node_pool_nodes_metadata) structure is documented below.

* `spec` - The specification of the node.  
  The [spec](#v2modelarts_node_pool_nodes_spec) structure is documented below.

* `status` - The status information of the node.  
  The [status](#v2modelarts_node_pool_nodes_status) structure is documented below.

<a name="v2modelarts_node_pool_nodes_metadata"></a>
The `metadata` block supports:

* `name` - The name of the node.

* `creation_timestamp` - The creation timestamp of the node.

* `labels` - The labels of the node.

* `annotations` - The annotation configuration of the node.

<a name="v2modelarts_node_pool_nodes_spec"></a>
The `spec` block supports:

* `flavor` - The flavor of the node.

* `extend_params` - The extend parameters of the node.

* `host_network` - The network configuration of the node.  
  The [host_network](#v2modelarts_node_pool_nodes_spec_host_network) structure is documented below.

* `os` - The OS information of the node.  
  The [os](#v2modelarts_node_pool_nodes_spec_os) structure is documented below.

<a name="v2modelarts_node_pool_nodes_spec_host_network"></a>
The `host_network` block supports:

* `vpc` - The VPC ID to which the node belongs.

* `subnet` - The subnet ID to which the node belongs.

* `security_groups` - The security group IDs that the node used.

<a name="v2modelarts_node_pool_nodes_spec_os"></a>
The `os` block supports:

* `name` - The OS name of the node.

* `image_id` - The image ID of the OS.

* `image_type` - The image type of the OS.

<a name="v2modelarts_node_pool_nodes_status"></a>
The `status` block supports:

* `phase` - The current phase of the node.
  + **Available**
  + **Creating**
  + **Deleting**
  + **Abnormal**
  + **Checking**

* `az` - The availability zone where the node is located.

* `driver` - The driver configuration of the node.  
  The [driver](#v2modelarts_node_pool_nodes_status_driver) structure is documented below.

* `os` - The OS information of the kubernetes node.  
  The [os](#v2modelarts_node_pool_nodes_status_os) structure is documented below.

* `plugins` - The plugin configuration of the node.  
  The [plugins](#v2modelarts_node_pool_nodes_status_plugins) structure is documented above.

* `private_ip` - The private IP address of the node.

* `resources` - The resource detail of the node, in JSON format.

* `available_resources` - The available resource detail of the node, in JSON format.

<a name="v2modelarts_node_pool_nodes_status_driver"></a>
The `driver` block supports:

* `phase` - The current phase of the driver.

* `version` - The version of the driver.

<a name="v2modelarts_node_pool_nodes_status_os"></a>
The `os` block supports:

* `name` - The OS name of the kubernetes node.

<a name="v2modelarts_node_pool_nodes_status_plugins"></a>
The `plugins` block supports:

* `name` - The name of the plugin.

* `phase` - The current phase of the plugin.

* `version` - The version of the plugin.
