---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_security\_group\_rule

Manages a security group rule resource within HuaweiCloud IEC.
This is an alternative to `huaweicloud_iec_security_group_rule_v1`

## Example Usage

```hcl
resource "huaweicloud_iec_security_group" "secgroup_test" {
  name        = "my_secgroup_test"
  description = "My security group"
}

resource "huaweicloud_iec_security_group_rule" "secgroup_rule_test" {
  direction         = "ingress"
  port_range_min    = 22
  port_range_max    = 22
  ethertype         = "IPv4"
  protocol          = "tcp"
  security_group_id = huaweicloud_iec_security_group.secgroup_test.id
  remote_ip_prefix  = "0.0.0.0/0"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String, ForceNew) Specifies a description for the 
    security group rule. Changing this parameter creates a new security group 
    rule resource.

* `direction` - (Required, String, ForceNew) The direction of the rule, valid 
    values are __ingress__ or __egress__. Changing this parameter creates a new 
    security group rule resource.

* `port_range_min` - (Optional, Int, ForceNew) The lower part of the allowed 
    port range, valid integer value needs to be between 1 and 65535. Changing 
    this parameter creates a new security group rule resource.

* `port_range_max` - (Optional, Int, ForceNew) The higher part of the allowed 
    port range, valid integer value needs to be between 1 and 65535. Changing 
    this parameter creates a new security group rule resource.

* `ethertype` - (Required, String, ForceNew) The layer 3 protocol type, valid 
    values are __IPv4__(IPv4 is default) or __IPv6__. Changing this parameter 
    creates a new security group rule resource.

* `protocol` - (Optional, String, ForceNew) The layer 4 protocol type, valid 
    values are following. Changing this parameter creates a new security group 
    rule resource. This is required if you want to specify a port range.
  * __tcp__
  * __udp__
  * __icmp__
  * __gre__

* `security_group_id` - (Required, String, ForceNew) The security group id the 
    rule should belong to, the value needs to be an Openstack ID of a security 
    group. Changing this parameter creates a new security group rule resource.

* `remote_ip_prefix` - (Optional, String, ForceNew) The remote CIDR, the value 
    needs to be a valid CIDR (i.e. 192.168.0.0/16).Changing this parameter 
    creates a new security group rule resource.

* `remote_group_id` - (Optional, String, ForceNew) The remote group id, the 
    value needs to be an Openstack ID of a security group. Changing this 
    parameter creates a new security group rule resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

IEC Security Group Rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_security_group_rule.secgroup_rule_test 2a02d1d3-437c-11eb-b721-fa163e8ac569
```
