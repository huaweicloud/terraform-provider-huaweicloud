---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_network_instances"
description: ""
---

# huaweicloud_cc_network_instances

Use this data source to get the list of CC network instances.

## Example Usage

```hcl
variable "network_instance_id" {} 

data "huaweicloud_cc_network_instances" "test" {
  network_instance_id = var.network_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `network_instance_id` - (Optional, String) Specifies the network instance ID.

* `name` - (Optional, String) Specifies the network instance name.

* `description` - (Optional, String) Specifies the network instance description.

* `cloud_connection_id` - (Optional, String) Specifies the cloud connection ID.

* `status` - (Optional, String) Specifies the status of the network instance.
  The options are as follows:
  + **ACTIVE**: The network instance is available.
  + **PENDING**: The network instance is being processed.
  + **ERROR**: The processing failed.

* `type` - (Optional, String) Specifies the type of the network instance.
  Value options are as follows:
  + **vpc**: a VPC.
  + **vgw**: a virtual gateway.

* `instance_id` - (Optional, String) Specifies the ID of the VPC or virtual gateway to be loaded to the cloud connection.

* `region_id` - (Optional, String) Specifies the region ID of the network instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `network_instances` - Network instance list.

  The [network_instances](#network_instances_struct) structure is documented below.

<a name="network_instances_struct"></a>
The `network_instances` block supports:

* `id` - Network instance ID.

* `name` - Network instance name.

* `description` - Network instance description.

* `domain_id` - ID of the account that the network instance belongs to.

* `created_at` - Time when the network instance was created.

* `updated_at` - Time when the network instance was updated.

* `cloud_connection_id` - The cloud connection ID.

* `instance_id` - The ID of the VPC or virtual gateway to be loaded to the cloud connection.

* `instance_domain_id` - Account ID of the VPC or virtual gateway.

* `region_id` - Region ID of the network instance.

* `project_id` - Project ID of the network instance.

* `status` - Status of the network instance.

* `type` - Type of the network instance.

* `cidrs` - The list of routes advertised by the network instance.
