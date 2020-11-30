---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_secgroup

Manages a Security Group resource within HuaweiCloud.
This is an alternative to `huaweicloud_networking_secgroup_v2`

## Example Usage

```hcl
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "secgroup_1"
  description = "My security group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the security group resource. If omitted, the provider-level region will be used. Changing this creates a new security group resource.

* `name` - (Required, String) A unique name for the security group.

* `description` - (Optional, String) Description for the security group.

* `tenant_id` - (Optional, String, ForceNew) The owner of the security group. Required if admin
    wants to create a port for another tenant. Changing this creates a new
    security group.

* `delete_default_rules` - (Optional, Bool, ForceNew) Whether or not to delete the default
    egress security rules. This is `false` by default. See the below note
    for more information.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the security group. Changing this creates a new security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Default Security Group Rules

In most cases, HuaweiCloud will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "huaweicloud_networking_secgroup_rule_v2" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
}

resource "huaweicloud_networking_secgroup_rule_v2" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
}
```

Please note that this behavior may differ depending on the configuration of
the HuaweiCloud cloud. The above illustrates the current default Neutron
behavior. Some HuaweiCloud clouds might provide additional rules and some might
not provide any rules at all (in which case the `delete_default_rules` setting
is moot).

## Timeouts
This resource provides the following timeouts configuration options:
- `delete` - Default is 10 minute.

## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
