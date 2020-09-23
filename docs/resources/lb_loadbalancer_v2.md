---
subcategory: "Elastic Load Balance (ELB)"
---

# huaweicloud\_lb\_loadbalancer\_v2

Manages a V2 loadbalancer resource within HuaweiCloud.

## Example Usage

### Basic Loadbalancer

```hcl
resource "huaweicloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```

### Loadbalancer With EIP

```hcl
resource "huaweicloud_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}

resource "huaweicloud_networking_eip_associate" "eip_1" {
  public_ip = "1.2.3.4"
  port_id   = huaweicloud_lb_loadbalancer_v2.lb_1.vip_port_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Human-readable name for the Loadbalancer. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the Loadbalancer.

* `vip_subnet_id` - (Required) The network on which to allocate the
    Loadbalancer's address. A tenant can only create Loadbalancers on networks
    authorized by policy (e.g. networks that belong to them or networks that
    are shared).  Changing this creates a new loadbalancer.

* `vip_address` - (Optional) The ip address of the load balancer.
    Changing this creates a new loadbalancer.

* `admin_state_up` - (Optional) The administrative state of the Loadbalancer.
    A valid value is true (UP) or false (DOWN).

## Attributes Reference

The following attributes are exported:

* `vip_subnet_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `vip_port_id` - The Port ID of the Load Balancer IP.
