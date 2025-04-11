---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_pools"
description: |-
  Use this data source to get the list of ELB pools.
---

# huaweicloud_lb_pools

Use this data source to get the list of ELB pools.

## Example Usage

```hcl
variable "pool_name" {}

data "huaweicloud_lb_pools" "test" {
  name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the ELB pool.

* `pool_id` - (Optional, String) Specifies the ID of the ELB pool.

* `description` - (Optional, String) Specifies the description of the ELB pool.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `loadbalancer_id` - (Optional, String) Specifies the loadbalancer ID of the ELB pool.

* `member_address` - (Optional, String)  Specifies the private IP address bound to the backend server.

* `member_device_id` - (Optional, String) Specifies the ID of the cloud server that serves as a backend server.

* `healthmonitor_id` - (Optional, String) Specifies the health monitor ID of the ELB pool.

* `protocol` - (Optional, String) Specifies the protocol of the ELB pool. This can either be TCP, UDP or HTTP.

* `lb_method` - (Optional, String) Specifies the method of the ELB pool. Must be one of ROUND_ROBIN, LEAST_CONNECTIONS,
  or SOURCE_IP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - The pool list. For details, see data structure of the pool field.
The [object](#pools_object) structure is documented below.

<a name="pools_object"></a>
The `pools` block supports:

* `id` - The pool ID.

* `name` - The pool name.

* `description` - The description of pool.

* `protocol` - The protocol of pool.

* `lb_method` - The load balancing algorithm to distribute traffic to the pool's members.

* `healthmonitor_id` - The health monitor ID of the LB pool.

* `protection_status` - Whether modification protection is enabled.

* `protection_reason` - The reason to enable modification protection.

* `listeners` - The listener list. The [object](#elem_object) structure is documented below.

* `loadbalancers` - The loadbalancer list. The [object](#elem_object) structure is documented below.

* `members` - The member list. The [object](#elem_object) structure is documented below.

* `persistence` - Indicates whether connections in the same session will be processed by the same pool member or not.
  The [object](#persistence_object) structure is documented below.

<a name="elem_object"></a>
The `listeners`,  `loadbalancers` or `members` block supports:

* `id` - The listener, loadbalancer or member ID.

<a name="persistence_object"></a>
The `persistence` block supports:

* `type` - The type of persistence mode.

* `cookie_name` - The name of the cookie if persistence mode is set appropriately.

* `timeout` - The sticky session timeout duration in minutes.
