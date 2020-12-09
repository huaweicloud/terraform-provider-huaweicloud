---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_secgroup\_rule

Manages a Security Group Rule resource within HuaweiCloud.
This is an alternative to `huaweicloud_networking_secgroup_rule_v2`

## Example Usage

```hcl
resource "huaweicloud_networking_secgroup" "mysecgroup" {
  name        = "secgroup"
  description = "My security group"
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.mysecgroup.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the security group rule resource. If omitted, the provider-level region will be used. Changing this creates a new security group rule resource.

* `direction` - (Required, String, ForceNew) The direction of the rule, valid values are __ingress__
    or __egress__. Changing this creates a new security group rule.

* `ethertype` - (Required, String, ForceNew) The layer 3 protocol type, valid values are __IPv4__
    or __IPv6__. Changing this creates a new security group rule.

* `protocol` - (Optional, String, ForceNew) The layer 4 protocol type, valid values are following. Changing this creates a new security group rule. This is required if you want to specify a port range.
  * __tcp__
  * __udp__
  * __icmp__
  * __ah__
  * __dccp__
  * __egp__
  * __esp__
  * __gre__
  * __igmp__
  * __ipv6-encap__
  * __ipv6-frag__
  * __ipv6-icmp__
  * __ipv6-nonxt__
  * __ipv6-opts__
  * __ipv6-route__
  * __ospf__
  * __pgm__
  * __rsvp__
  * __sctp__
  * __udplite__
  * __vrrp__

* `port_range_min` - (Optional, String, ForceNew) The lower part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `port_range_max` - (Optional, Int, ForceNew) The higher part of the allowed port range, valid
    integer value needs to be between 1 and 65535. Changing this creates a new
    security group rule.

* `remote_ip_prefix` - (Optional, String, ForceNew) The remote CIDR, the value needs to be a valid
    CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule.

* `remote_group_id` - (Optional, String, ForceNew) The remote group id, the value needs to be an
    Openstack ID of a security group in the same tenant. Changing this creates
    a new security group rule.

* `security_group_id` - (Required, String, ForceNew) The security group id the rule should belong
    to, the value needs to be an Openstack ID of a security group in the same
    tenant. Changing this creates a new security group rule.

* `tenant_id` - (Optional, String, ForceNew) The owner of the security group. Required if admin
    wants to create a port for another tenant. Changing this creates a new
    security group rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Timeouts
This resource provides the following timeouts configuration options:
- `delete` - Default is 10 minute.

## Import

Security Group Rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
