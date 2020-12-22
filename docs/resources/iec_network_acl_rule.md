---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_network\_acl\_rule

Manages a network ACL rule resource within HuaweiCloud IEC.

## Example Usage

```hcl
resource "huaweicloud_iec_network_acl_rule" "rule_test" {
  protocol               = "tcp"
  action                 = "deny"
  source_ip_address      = "112.25.96.0/20"
  source_port            = "445"
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, String) Specifies a unique id for the iec 
    network ACL.

* `direction` - (Required, String, ForceNew) Specifies the direction of the 
    rule, valid values are *ingress* or *egress*. Changing this parameter 
    creates a new iec network ACL rule resource.

* `description` - (Optional, String) Specifies the description for the iec 
    network ACL rule.

* `protocol` - (Optional, String) Specifies the protocol supported by the iec 
    network ACL rule. Valid values are: *tcp*, *udp*, *icmp* and *any*.

* `action` - (Optional, String) Specifies the action in the iec network ACL 
    rule. Currently, the value can be *allow* or *deny*.

* `ip_version` - (Optional, Int) Specifies the IP version, either 4 (default) 
    or 6. This parameter is available after the IPv6 function is enabled.

* `source_ip_address` - (Optional, String) Specifies the source IP address that 
    the traffic is allowed from. The default value is *0.0.0.0/0*. For example: 
    xxx.xxx.xxx.xxx (IP address), xxx.xxx.xxx.0/24 (CIDR block).

* `destination_ip_address` - (Optional, String) Specifies the destination IP 
    address to which the traffic is allowed. The default value is *0.0.0.0/0*. 
    For example: xxx.xxx.xxx.xxx (IP address), xxx.xxx.xxx.0/24 (CIDR block).

* `source_port` - (Optional, String) Specifies the source port number or port   
    number range. The value ranges from 1 to 65535. For a port number range, 
    enter two port numbers connected by a hyphen (-). For example, 1-100.

* `destination_port` - (Optional, String) Specifies the destination port number 
    or port number range. The value ranges from 1 to 65535. For a port number 
    range, enter two port numbers connected by a hyphen (-). For example, 1-100.

* `enabled` - (Optional, Bool) Enabled status for the iec network ACL rule. 
    Defaults to true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `poilicy_id` - Specifies the ID of the firewall policy for the iec network 
    ACL.

## Import

network ACL rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_network_acl_rule.rule_test 89a84b28-4cc2-4859-9885-c67e802a46a3
```
