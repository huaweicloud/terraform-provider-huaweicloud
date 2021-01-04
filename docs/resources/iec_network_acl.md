---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_network\_acl

Manages a network ACL resource within HuaweiCloud IEC.

## Example Usage

### Without networks

```hcl
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_network_acl" "acl_test" {
  name        = "acl_demo"
  description = "This is a network ACL of IEC without networks."
}
```

### With networks

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
}

resource "huaweicloud_iec_network_acl" "acl_test" {
  name        = "acl_demo"
  description = "This is a network ACL of IEC with networks."
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

* `networks` - (Optional, List) Specifies an list of one or more networks. 
    The networks object structure is documented below.
    
The `networks` block supports:

* `vpc_id` - (Required, String) Specifies the id of the iec vpc.

* `subnet_id` - (Required, String) Specifies the id of the iec subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the iec network ACL. 

* `inbound_rules` - A list of the IDs of ingress rules associated with the 
    iec network ACL.

* `outbound_rules` - A list of the IDs of egress rules associated with the 
    iec network ACL.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.
