---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_servers"
description: |-
  Use this data source to get servers of CPH phone.
---

# huaweicloud_cph_servers

Use this data source to get servers of CPH phone.

## Example Usage

```hcl
data "huaweicloud_cph_servers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_name` - (Optional, String) Specifies the cloud phone server name and support fuzzy query.

* `server_id` - (Optional, String) Specifies the cloud phone server ID.

* `network_version` - (Optional, String) Specifies whether the cloud phone server is a custom network identifier.
  + **v1**: System-defined network cloud phone server.
  + **v2**: Cloud phone server for custom network.

* `phone_flavor` - (Optional, String) Specifies the cloud phone flavor name.

* `create_since` - (Optional, String) The creation start time. For example, **2024-10-15 15:04:05**.

* `create_until` - (Optional, String) The creation end time. For example, **2024-10-15 15:04:05**.

* `status` - (Optional, String) Specifies the server status.
  + **0, 1, 3, 4**: Creating
  + **2**: Abnormal
  + **5**: Normal
  + **8**: Freeze
  + **10**: Shut down
  + **11**: Shutting down
  + **12**: Shutdown failed
  + **13**: Starting up

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - The cloud phone server list.

  The [servers](#servers_struct) structure is documented below.

<a name="servers_struct"></a>
The `servers` block supports:

* `server_id` - The cloud phone server ID.

* `server_flavor` - The cloud phone server flavor name.

* `keypair_name` - The name of the key pair used to connect to the cloud phone.

* `subnet_id` - The ID of the subnet to which the cloud phone server belongs.

* `subnet_cidr` - The subnet CIDR to which the cloud phone server belongs.

* `addresses` - The IP related information of the cloud phone server.

  The [addresses](#servers_addresses_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID to which the cloud phone server belongs.

* `availability_zone` - The availability zone where the cloud mobile server is located.

* `phone_flavor` - The cloud phone flavor name.

* `cidr` - The network segment of VPC to which the cloud phone server belongs.

* `network_version` - Whether the cloud phone server is a custom network identifier.

* `status` - The cloud phone server status.

* `vpc_cidr` - The VPC CIDR.
  When the value of `network_version` is **v1**, it indicates the VPC CIDR of the resource tenant to which
  the cloud mobile server belongs; when the value of `network_version` is **v2**, it indicates the VPC CIDR of the VPC
  specified by the tenant when creating the server.

* `resource_project_id` - The project ID of the cloud phone server.

* `metadata` - The order and product related information.

  The [metadata](#servers_metadata_struct) structure is documented below.

* `update_time` - The update time.

* `server_name` - The cloud phone server name.

* `vpc_id` - The ID of the virtual private cloud (VPC for short) to which the cloud mobile server belongs.
  When the value of `network_version` is **v1**, it indicates the VPC ID of the resource tenant to which
  the cloud mobile server belongs; when the value of `network_version` is **v2**, it indicates the VPC ID
  of the VPC specified by the tenant when creating the server.

* `create_time` - The creation time.

<a name="servers_addresses_struct"></a>
The `addresses` block supports:

* `server_ip` - The intranet IP of cloud phone server.

* `public_ip` - The public IP of cloud phone server.

<a name="servers_metadata_struct"></a>
The `metadata` block supports:

* `product_id` - The product ID.

* `order_id` - The order ID.
