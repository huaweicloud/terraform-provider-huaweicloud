---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_network\_acl

Manages a network ACL resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc_demo"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name        = "subnet_demo"
  cidr        = "192.168.128.0/18"
  vpc_id      = huaweicloud_iec_vpc.vpc_test.id
  site_id     = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip  = "192.168.128.3"
}

resource "huaweicloud_iec_network_acl_rule" "rule_1" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "ingress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "132.156.0.0/16"
  destination_ip_address = "192.168.128.0/18"
  destination_port       = "445"
  enabled                = true
}

resource "huaweicloud_iec_network_acl_rule" "rule_2" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "egress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "192.168.128.0/18"
  destination_ip_address = "152.16.30.0/24"
  destination_port       = "45"
  enabled                = true
}

resource "huaweicloud_iec_network_acl" "acl_test" {
  name = "acl_demo"
  networks {
    vpc_id = huaweicloud_iec_vpc.vpc_test.id
    subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the iec network ACL name. This 
    parameter can contain a maximum of 64 characters, which may consist of 
    letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional, String) Specifies the supplementary information 
    about the iec network ACL. This parameter can contain a maximum of 255 
    characters and cannot contain angle brackets (< or >).

* `inbound_rules` - (Optional, List)  A list of the IDs of ingress rules 
    associated with the iec network ACL. The maximum length of the list is 10.

* `outbound_rules` - (Optional, List) A list of the IDs of egress rules 
    associated with the iec network ACL. The maximum length of the list is 10

* `networks` - (Optional, Set) An Set of one or more networks. The networks 
    object structure is documented below.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the iec network ACL. 

The `networks` block supports:

* `vpc_id` - The id of the iec vpc.
* `subnet_id` - The id of the iec subnet.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.
