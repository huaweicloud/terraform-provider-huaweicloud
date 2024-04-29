---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_members"
description: ""
---

# huaweicloud_elb_members

Use this data source to get the list of ELB members.

## Example Usage

```hcl
variable "member_name" {}
variable "elb_pool_id" {}

data "huaweicloud_elb_members" "test" {
  pool_id = var.elb_pool_id
  name    = var.member_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `pool_id` - (Required, String) Specifies the ID of the backend server group.

* `subnet_id` - (Optional, String) Specifies the ID of the IPv4 or IPv6 subnet where the backend server resides.

* `member_id` - (Optional, String) Specifies the backend server ID.

* `name` - (Optional, String) Specifies the backend server name.

* `address` - (Optional, String) Specifies the IP address bound to the backend server.

* `protocol_port` - (Optional, Int) Specifies the port used by the backend server to receive requests.

* `weight` - (Optional, Int)  Specifies the weight of the backend server. Requests are routed to backend servers in the
  same backend server group based on their weights.

* `member_type` - (Optional, String) Specifies the type of the backend server. The valid values are as follows:
  + **ip**: IP as backend servers.
  + **instance**: ECSs used as backend servers Multiple values can be queried in the format of
     member_type=xxx&member_type=xxx.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `members` - Lists the members.
  The [members](#Elb_members) structure is documented below.

<a name="Elb_members"></a>
The `members` block supports:

* `name` - The backend server name.

* `id` - The backend server ID.

* `subnet_id` - The ID of the IPv4 or IPv6 subnet where the backend server resides.

* `protocol_port` - The port used by the backend server to receive requests.

* `weight` - The weight of the backend server. Requests are routed to backend servers in the same backend server group
  based on their weights.

* `address` - The private IP address bound to the backend server.

* `member_type` - The type of the backend server.
