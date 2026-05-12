---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_peering_connections"
description: |-
  Use this data source to get the list of VPC peering connections.
---

# huaweicloud_vpc_peering_connections

Use this data source to get the list of VPC peering connections in the current tenant.

~> Currently, the maximum of 2000 peering connections can be queried.

## Example Usage

### Query all peering connections and without any filter

```hcl
data "huaweicloud_vpc_peering_connections" "all" {}
```

### Query the peering connections and using name filter

```hcl
variable "peering_connection_name" {}

data "huaweicloud_vpc_peering_connections" "by_vpc" {
  name = var.peering_connection_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the peering connections are located.  
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the ID of the VPC peering connection to be queried.

* `name` - (Optional, String) Specifies the name of the VPC peering connection to be queried.

* `status` - (Optional, String) Specifies the status of the VPC peering connection to be queried.  
  The valid values are as follows:
  + **PENDING_ACCEPTANCE**
  + **REJECTED**
  + **EXPIRED**
  + **DELETED**
  + **ACTIVE**

* `project_id` - (Optional, String) Specifies the project ID of the VPC to be queried.

* `vpc_id` - (Optional, String) Specifies the ID of the requester's VPC to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of peering connections that matched the filter parameters.  
  The [connections](#vpc_peering_connections_attr) structure is documented below.

<a name="vpc_peering_connections_attr"></a>
The `connections` block supports:

* `id` - The ID of the peering connection.

* `name` - The name of the peering connection.

* `status` - The status of the peering connection.

* `description` - The description of the peering connection.

* `request_vpc_info` - The information of the requester's VPC.  
  The [request_vpc_info](#vpc_peering_connections_request_vpc_info_attr) structure is documented below.

* `accept_vpc_info` - The information of the accepter's VPC.  
  The [accept_vpc_info](#vpc_peering_connections_accept_vpc_info_attr) structure is documented below.

* `created_at` - The creation time of the peering connection, in RFC3339 format.

* `updated_at` - The latest update time of the peering connection, in RFC3339 format.

<a name="vpc_peering_connections_request_vpc_info_attr"></a>
The `request_vpc_info` block supports:

* `vpc_id` - The ID of the requester's VPC.

* `project_id` - The project ID of the requester to which the VPC of peering connection belongs.

<a name="vpc_peering_connections_accept_vpc_info_attr"></a>
The `accept_vpc_info` block supports:

* `vpc_id` - The ID of the accepter's VPC.

* `project_id` - The project ID of the accepter to which the VPC of peering connection belongs.
