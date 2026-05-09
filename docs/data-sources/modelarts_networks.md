---
subcategory: "ModelArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_networks"
description: |-
  Use this data source to query the ModelArts networks within HuaweiCloud.
---

# huaweicloud_modelarts_networks

Use this data source to query the ModelArts networks within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_modelarts_networks" "test" {}
```

### Filter by label selector

```hcl
variable "label_selector" {}

data "huaweicloud_modelarts_networks" "test" {
  label_selector = var.label_selector
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the networks are located.  
  If omitted, the provider-level region will be used.

* `label_selector` - (Optional, String) Specifies the label selector to filter networks.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `networks` - The list of networks that matched filter parameters.  
  The [networks](#networks_attr) structure is documented below.

<a name="networks_attr"></a>
The `networks` block supports:

* `api_version` - The API version of the network.

* `kind` - The kind of the network.

* `metadata` - The metadata of the network.  
  The [metadata](#networks_metadata_attr) structure is documented below.

* `spec` - The spec of the network.  
  The [spec](#networks_spec_attr) structure is documented below.

* `status` - The status of the network.  
  The [status](#networks_status_attr) structure is documented below.

<a name="networks_metadata_attr"></a>
The `metadata` block supports:

* `name` - The name of the network.

* `created_at` - The creation time of the network.

* `labels` - The labels of the network.  
  The [labels](#networks_metadata_labels_attr) structure is documented below.

* `annotations` - The annotations of the network.  
  The [annotations](#networks_metadata_annotations_attr) structure is documented below.

<a name="networks_metadata_labels_attr"></a>
The `labels` block supports:

* `os_modelarts_name` - The display name of the resource pool.

* `os_modelarts_workspace_id` - The workspace ID of the resource pool.

<a name="networks_metadata_annotations_attr"></a>
The `annotations` block supports:

* `os_modelarts_description` - The description of the network.

<a name="networks_spec_attr"></a>
The `spec` block supports:

* `cidr` - The CIDR block of the network.

* `connection` - The connection information of the network.  
  The [connection](#networks_spec_connection_attr) structure is documented below.

<a name="networks_spec_connection_attr"></a>
The `connection` block supports:

* `peer_connection_list` - The peer connection list of the network.  
  The [peer_connection_list](#networks_spec_connection_peer_connection_item_attr) structure is documented below.

<a name="networks_spec_connection_peer_connection_item_attr"></a>
The `peer_connection_list` block supports:

* `peer_vpc_id` - The ID of the peer VPC.

* `peer_subnet_id` - The ID of the peer subnet.

* `default_gateway` - Whether to create the default gateway.

<a name="networks_status_attr"></a>
The `status` block supports:

* `phase` - The phase of the network.  
  The valid values are as follows:
  + **Creating**
  + **Active**
  + **Abnormal**

* `connection_status` - The connection status of the network.  
  The [connection_status](#networks_status_connection_status_attr) structure is documented below.

<a name="networks_status_connection_status_attr"></a>
The `connection_status` block supports:

* `peer_connection_status` - The peer connection status list of the network.  
  The [peer_connection_status](#networks_status_peer_connection_status_attr) structure is documented below.

* `sfs_turbo_status` - The SFS Turbo connection status list of the network.  
  The [sfs_turbo_status](#networks_status_sfs_turbo_status_attr) structure is documented below.

<a name="networks_status_peer_connection_status_attr"></a>
The `peer_connection_status` block supports:

* `peer_vpc_id` - The ID of the peer VPC.

* `peer_subnet_id` - The ID of the peer subnet.

* `default_gateway` - Whether the default gateway is created.

* `phase` - The connection phase of the peer connection.  
  The valid values are as follows:
  + **Connecting**
  + **Active**
  + **Abnormal**

<a name="networks_status_sfs_turbo_status_attr"></a>
The `sfs_turbo_status` block supports:

* `name` - The name of the SFS Turbo instance.

* `sfs_id` - The ID of the SFS Turbo instance.

* `connection_type` - The connection type of the SFS Turbo.  
  The valid values are as follows:
  + **VpcPort**
  + **Peering**

* `ip_addr` - The IP address of the SFS Turbo.

* `status` - The connection status of the SFS Turbo.  
  The valid values are as follows:
  + **Active**
  + **Abnormal**
  + **Creating**
  + **Deleting**
