---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_secgroup"
description: ""
---

# huaweicloud_networking_secgroup

Manages a Security Group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "secgroup_1"
  description = "My security group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the security group resource. If omitted, the
  provider-level region will be used. Changing this creates a new security group resource.

* `name` - (Required, String) Specifies a unique name for the security group.

* `description` - (Optional, String) Specifies the description for the security group.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the security group.
  Changing this creates a new security group.

* `delete_default_rules` - (Optional, Bool, ForceNew) Specifies whether or not to delete the default security rules.
  This is `false` by default.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the security group.

-> **NOTE:** The default security rules are described
in [HuaweiCloud](https://support.huaweicloud.com/intl/en-us/usermanual-vpc/SecurityGroup_0003.html). See the below
section for more information.

## Default Security Group Rules

In most cases, HuaweiCloud will create some security group rules for each new security group. These security group rules
will not be managed by Terraform, so if you prefer to have *all*
aspects of your infrastructure managed by Terraform, set `delete_default_rules` to `true`
and then create separate security group rules such as the following:

```hcl
resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_v4" {
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  direction         = "egress"
  ethertype         = "IPv4"
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_v6" {
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  direction         = "egress"
  ethertype         = "IPv6"
}

resource "huaweicloud_networking_secgroup_rule" "allow_ssh" {
  security_group_id = huaweicloud_networking_secgroup.secgroup.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `rules` - The array of security group rules associating with the security group.
  The [rule object](#security_group_rule) is documented below.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.

<a name="security_group_rule"></a>
The `rules` block supports:

* `id` - The security group rule ID.
* `description` - The supplementary information about the security group rule.
* `direction` - The direction of the rule. The value can be *egress* or *ingress*.
* `ethertype` - The IP protocol version. The value can be *IPv4* or *IPv6*.
* `protocol` - The protocol type.
* `ports` - The port value range.
* `remote_ip_prefix` - The remote IP address. The value can be in the CIDR format or IP addresses.
* `remote_group_id` - The ID of the peer security group.
* `remote_address_group_id` - The ID of the remote address group.
* `action` - The effective policy.
* `priority` - The priority number.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minutes.

## Import

Security Groups can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
