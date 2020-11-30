---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_port

Use this data source to get the ID of an available HuaweiCloud port.
This is an alternative to `huaweicloud_networking_port_v2`

## Example Usage

```hcl
data "huaweicloud_networking_port" "port_1" {
  name = "port_1"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve port ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional, String) The owner of the port.

* `port_id` - (Optional, String) The ID of the port.

* `name` - (Optional, String) The name of the port.

* `admin_state_up` - (Optional, Bool) The administrative state of the port.

* `network_id` - (Optional, String) The ID of the network the port belongs to.

* `device_owner` - (Optional, String) The device owner of the port.

* `mac_address` - (Optional, String) The MAC address of the port.

* `device_id` - (Optional, String) The ID of the device the port belongs to.

* `fixed_ip` - (Optional, String) The port IP address filter.

* `status` - (Optional, String) The status of the port.

* `security_group_ids` - (Optional, String) The list of port security group IDs to filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.

* `all_security_group_ids` - The set of security group IDs applied on the port.
