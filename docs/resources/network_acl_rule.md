---
subcategory: "Network ACL"
---

# huaweicloud\_network\_acl\_rule

Manages a network ACL rule resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_network_acl_rule" "rule_1" {
  name                   = "rule_1"
  protocol               = "udp"
  action                 = "deny"
  source_ip_address      = "1.2.3.4"
  source_port            = "444"
  destination_ip_address = "4.3.2.0/24"
  destination_port       = "555"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the network ACL rule resource. If omitted, the provider-level region will work as default. Changing this creates a new network ACL rule resource.

* `name` - (Optional) Specifies a unique name for the network ACL rule.

* `description` - (Optional) Specifies the description for the network ACL rule.

* `protocol` - (Required) Specifies the protocol supported by the network ACL rule.
     Valid values are: *tcp*, *udp*, *icmp* and *any*.

* `action` - (Required) Specifies the action in the network ACL rule. Currently, the value can be *allow* or *deny*.

* `ip_version` - (Optional) Specifies the IP version, either 4 (default) or 6. This parameter is
    available after the IPv6 function is enabled.

* `source_ip_address` - (Optional) Specifies the source IP address that the traffic is allowed from.
    The default value is *0.0.0.0/0*. For example: xxx.xxx.xxx.xxx (IP address), xxx.xxx.xxx.0/24 (CIDR block).

* `destination_ip_address` - (Optional) Specifies the destination IP address to which the traffic is allowed.
    The default value is *0.0.0.0/0*. For example: xxx.xxx.xxx.xxx (IP address), xxx.xxx.xxx.0/24 (CIDR block).

* `source_port` - (Optional) Specifies the source port number or port number range. The value ranges from 1 to 65535.
    For a port number range, enter two port numbers connected by a hyphen (-). For example, 1-100.

* `destination_port` - (Optional) Specifies the destination port number or port number range. The value ranges from 1 to 65535.
    For a port number range, enter two port numbers connected by a hyphen (-). For example, 1-100.

* `enabled` - (Optional) Enabled status for the network ACL rule. Defaults to true.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `action` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `source_ip_address` - See Argument Reference above.
* `destination_ip_address` - See Argument Reference above.
* `source_port` - See Argument Reference above.
* `destination_port` - See Argument Reference above.
* `enabled` - See Argument Reference above.

## Import

network ACL rules can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_network_acl_rule.rule_1 89a84b28-4cc2-4859-9885-c67e802a46a3
```
