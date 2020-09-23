---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_port

Manages a Port resource within HuaweiCloud.
This is an alternative to `huaweicloud_networking_port_v2`

## Example Usage

```hcl
data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

resource "huaweicloud_networking_port" "myport" {
  name           = "port"
  network_id     = data.huaweicloud_vpc_subnet.mynet.id
  admin_state_up = "true"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique name for the port. Changing this
    updates the `name` of an existing port.

* `network_id` - (Required) The ID of the network to attach the port to. Changing
    this creates a new port.

* `admin_state_up` - (Optional) Administrative up/down status for the port
    (must be "true" or "false" if provided). Changing this updates the
    `admin_state_up` of an existing port.

* `mac_address` - (Optional) Specify a specific MAC address for the port. Changing
    this creates a new port.

* `tenant_id` - (Optional) The owner of the Port. Required if admin wants
    to create a port for another tenant. Changing this creates a new port.

* `device_owner` - (Optional) The device owner of the Port. Changing this creates
    a new port.

* `security_group_ids` - (Optional - Conflicts with `no_security_groups`) A list
    of security group IDs to apply to the port. The security groups must be
    specified by ID and not name (as opposed to how they are configured with
    the Compute Instance).

* `no_security_groups` - (Optional - Conflicts with `security_group_ids`) If set to
    `true`, then no security groups are applied to the port. If set to `false` and
    no `security_group_ids` are specified, then the Port will yield to the default
    behavior of the Networking service, which is to usually apply the "default"
    security group.

* `device_id` - (Optional) The ID of the device attached to the port. Changing this
    creates a new port.

* `fixed_ip` - (Optional) An array of desired IPs for this port. The structure is
    described below.

* `allowed_address_pairs` - (Optional) An IP/MAC Address pair of additional IP
    addresses that can be active on this port. The structure is described
    below.

* `extra_dhcp_option` - (Optional) An extra DHCP option that needs to be configured
    on the port. The structure is described below. Can be specified multiple
    times.

* `value_specs` - (Optional) Map of additional options.

The `fixed_ip` block supports:

* `subnet_id` - (Required) Subnet in which to allocate IP address for
this port.

* `ip_address` - (Optional) IP address desired in the subnet for this port. If
you don't specify `ip_address`, an available IP address from the specified
subnet will be allocated to this port. This field will not be populated if it
is left blank. To retrieve the assigned IP address, use the `all_fixed_ips`
attribute.

The `allowed_address_pairs` block supports:

* `ip_address` - (Required) The additional IP address.

* `mac_address` - (Optional) The additional MAC address.

The `extra_dhcp_option` block supports:

* `name` - (Required) Name of the DHCP option.

* `value` - (Required) Value of the DHCP option.

* `ip_version` - (Optional) IP protocol version. Defaults to 4.

## Attributes Reference

The following attributes are exported:

* `admin_state_up` - See Argument Reference above.
* `mac_address` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `device_owner` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `device_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.
* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.
* `all_security_group_ids` - The collection of Security Group IDs on the port
  which have been explicitly and implicitly added.
* `extra_dhcp_option` - See Argument Reference above.

## Import

Ports can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_networking_port.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```
