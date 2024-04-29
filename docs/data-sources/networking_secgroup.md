---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_secgroup"
description: ""
---

# huaweicloud_networking_secgroup

Use this data source to get the ID of an available HuaweiCloud security group.

## Example Usage

```hcl
data "huaweicloud_networking_secgroup" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the security group. If omitted, the
  provider-level region will be used.

* `secgroup_id` - (Optional, String) Specifies the ID of the security group.

* `name` - (Optional, String) Specifies the name of the security group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the security group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `description`- The description of the security group.

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
