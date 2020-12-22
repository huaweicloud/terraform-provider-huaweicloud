---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_network\_acl

Manages a network ACL resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "iec-vpc-test"
  cidr = "192.168.0.0/16"
  mode = "SYSTEM"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name        = "iec-subnet-test"
  cidr        = "192.168.199.0/24"
  vpc_id      = huaweicloud_iec_vpc.vpc_test.id
  site_id     = data.huaweicloud_iec_sites.iec_sites.sites[0].id
}

resource "huaweicloud_iec_network_acl_rule" "rule_test_1" {
  protocol               = "tcp"
  action                 = "deny"
  source_ip_address      = "112.25.96.0/20"
  source_port            = "445"
}

resource "huaweicloud_iec_network_acl_rule" "rule_test_2" {
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "124.70.64.0/20"
  source_port            = "8080"
}

resource "huaweicloud_iec_network_acl" "acl_test" {
  name          = "iec-acl-test"
  inbound_rules = [huaweicloud_iec_network_acl_rule.rule_test_1.id,
    huaweicloud_iec_network_acl_rule.rule_test_2.id]
  networks {
    vpc_id    = huaweicloud_iec_vpc.vpc_test.id
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
    associated with the iec network ACL. 

* `outbound_rules` - (Optional, List) A list of the IDs of egress rules 
    associated with the iec network ACL. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the network ACL.

* `networks` - An Set of one or more networks. The networks object structure is 
    documented below.

* `status` - The status of the iec network ACL. 

The `networks` block supports:

* `vpc_id` - The id of the iec vpc.
* `subnet_id` - The id of the iec subnet.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

IEC network ACL can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_network_acl.acl_test 773d965e-43fb-11eb-b721-fa163e8ac569
```
