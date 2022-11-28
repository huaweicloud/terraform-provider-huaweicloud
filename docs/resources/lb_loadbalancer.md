---
subcategory: "Elastic Load Balance (ELB)"
---

# huaweicloud_lb_loadbalancer

Manages an ELB loadbalancer resource within HuaweiCloud.

## Example Usage

### Basic Loadbalancer

```hcl
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_loadbalancer" "lb_1" {
  vip_subnet_id = var.ipv4_subnet_id

  tags = {
    key = "value"
  }
}
```

### Loadbalancer With EIP

```hcl
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_loadbalancer" "lb_1" {
  vip_subnet_id = ipv4_subnet_id
}

resource "huaweicloud_vpc_eip_associate" "eip_1" {
  public_ip = "1.2.3.4"
  port_id   = huaweicloud_lb_loadbalancer.lb_1.vip_port_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the loadbalancer resource. If omitted, the
  provider-level region will be used. Changing this creates a new loadbalancer.

* `name` - (Optional, String) Human-readable name for the loadbalancer. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the loadbalancer.

* `vip_subnet_id` - (Required, String, ForceNew) The **IPv4 subnet ID** of the subnet where the load balancer works.
  Changing this creates a new loadbalancer.

* `vip_address` - (Optional, String, ForceNew) The ip address of the load balancer. Changing this creates a new
  loadbalancer.

* `admin_state_up` - (Optional, Bool) The administrative state of the loadbalancer. A valid value is true (UP) or
  false (DOWN).

* `tags` - (Optional, Map) The key/value pairs to associate with the loadbalancer.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the loadbalancer. Changing this
  creates a new loadbalancer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vip_port_id` - The Port ID of the Load Balancer IP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 5 minute.

## Import

Load balancers can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_lb_loadbalancer.test 3e3632db-36c6-4b28-a92e-e72e6562daa6
```
