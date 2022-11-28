---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# huaweicloud_elb_loadbalancer

Manages a Dedicated Load Balancer resource within HuaweiCloud.

## Example Usage

### Basic Loadbalancer

```hcl
resource "huaweicloud_elb_loadbalancer" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id         = "{{ vpc_id }}"
  ipv4_subnet_id = "{{ ipv4_subnet_id }}"

  l4_flavor_id = "{{ l4_flavor_id }}"
  l7_flavor_id = "{{ l7_flavor_id }}"

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]

  enterprise_project_id = "{{ eps_id }}"
}
```

### Loadbalancer With Existing EIP

```hcl
resource "huaweicloud_elb_loadbalancer" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id            = "{{ vpc_id }}"
  ipv6_network_id   = "{{ ipv6_network_id }}"
  ipv6_bandwidth_id = "{{ ipv6_bandwidth_id }}"
  ipv4_subnet_id    = "{{ ipv4_subnet_id }}"

  l4_flavor_id = "{{ l4_flavor_id }}"
  l7_flavor_id = "{{ l7_flavor_id }}"

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]

  enterprise_project_id = "{{ eps_id }}"

  ipv4_eip_id = "{{ eip_id }}"
}
```

### Loadbalancer With EIP

```hcl
resource "huaweicloud_elb_loadbalancer" "basic" {
  name              = "basic"
  description       = "basic example"
  cross_vpc_backend = true

  vpc_id            = "{{ vpc_id }}"
  ipv6_network_id   = "{{ ipv6_network_id }}"
  ipv6_bandwidth_id = "{{ ipv6_bandwidth_id }}"
  ipv4_subnet_id    = "{{ ipv4_subnet_id }}"

  l4_flavor_id = "{{ l4_flavor_id }}"
  l7_flavor_id = "{{ l7_flavor_id }}"

  availability_zone = [
    "cn-north-4a",
    "cn-north-4b",
  ]

  enterprise_project_id = "{{ eps_id }}"

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new loadbalancer.

* `availability_zone` - (Required, List, ForceNew) Specifies the list of AZ names. Changing this parameter will create a
  new resource.

* `name` - (Required, String) Human-readable name for the loadbalancer.

* `description` - (Optional, String) Human-readable description for the loadbalancer.

* `cross_vpc_backend` - (Optional, Bool) Enable this if you want to associate the IP addresses of backend servers with
  your load balancer. Can only be true when updating.

* `vpc_id` - (Optional, String, ForceNew) The vpc on which to create the loadbalancer. Changing this creates a new
  loadbalancer.

* `ipv4_subnet_id` - (Optional, String) The **IPv4 subnet ID** of the subnet on which to allocate the loadbalancer's
  ipv4 address.

* `ipv6_network_id` - (Optional, String) The **ID** of the subnet on which to allocate the loadbalancer's ipv6 address.

* `ipv6_bandwidth_id` - (Optional, String) The ipv6 bandwidth id. Only support shared bandwidth.

* `ipv4_address` - (Optional, String) The ipv4 address of the load balancer.

* `ipv4_eip_id` - (Optional, String, ForceNew) The ID of the EIP. Changing this parameter will create a new resource.

-> **NOTE:** If the ipv4_eip_id parameter is configured, you do not need to configure the bandwidth parameters:
`iptype`, `bandwidth_charge_mode`, `bandwidth_size` and `share_type`.

* `iptype` - (Optional, String, ForceNew) Elastic IP type. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional, String, ForceNew) Bandwidth billing type. Changing this parameter will create a
  new resource.

* `sharetype` - (Optional, String, ForceNew) Bandwidth sharing type. Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional, Int, ForceNew) Bandwidth size. Changing this parameter will create a new resource.

* `l4_flavor_id` - (Optional, String) The L4 flavor id of the load balancer.

* `l7_flavor_id` - (Optional, String) The L7 flavor id of the load balancer.

* `tags` - (Optional, Map) The key/value pairs to associate with the loadbalancer.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the loadbalancer. Changing this
  creates a new loadbalancer.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the ELB loadbalancer.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.
  Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the ELB loadbalancer.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the ELB loadbalancer.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are **true** and **false**.

* `autoscaling_enabled` - (Optional, Bool) Specifies whether autoscaling is enabled. Valid values are **true** and
  **false**.

* `min_l7_flavor_id` - (Optional, String) Specifies the ID of the minimum Layer-7 flavor for elastic scaling.
  This parameter cannot be left blank if there are HTTP or HTTPS listeners.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ipv4_eip` - The ipv4 eip address of the Load Balancer.
* `ipv6_eip` - The ipv6 eip address of the Load Balancer.
* `ipv6_eip_id` - The ipv6 eip id of the Load Balancer.
* `ipv6_address` - The ipv6 address of the Load Balancer.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 5 minute.

## Import

ELB loadbalancer can be imported using the loadbalancer ID, e.g.

```
$ terraform import huaweicloud_elb_loadbalancer.loadbalancer_1 5c20fdad-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `ipv6_bandwidth_id`, `iptype`,
`bandwidth_charge_mode`, `sharetype` and `bandwidth_size`.
It is generally recommended running `terraform plan` after importing a loadbalancer.
You can then decide if changes should be applied to the loadbalancer, or the resource
definition should be updated to align with the loadbalancer. Also you can ignore changes as below.

```
resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
    ...
  lifecycle {
    ignore_changes = [
      ipv6_bandwidth_id, iptype, bandwidth_charge_mode, sharetype, bandwidth_size,
    ]
  }
}
```
