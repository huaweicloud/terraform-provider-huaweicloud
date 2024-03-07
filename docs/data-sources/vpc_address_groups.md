---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_address_groups

Use this data source to get a list of VPC IP address groups.

## Example Usage: Query All Address Groups

```hcl
data "huaweicloud_vpc_address_groups" "demo1" {
}
```

## Example Usage: Query An Address Group by Name

```hcl
data "huaweicloud_vpc_address_groups" "demo2" {
  name = "your_address_group_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resources.
  If omitted, the provider-level region will be used.

* `description` - (Optional, String) Provides supplementary information about an IP address group,
  which can be used to filter the IP address group.

* `group_id` - (Optional, String) Unique ID of an IP address group, which can be used to filter the IP address group.

* `ip_version` - (Optional, Int) Version of IP addresses in an IP address group,
  which can be used to filter the IP address group.

* `name` - (Optional, String) Name of an IP address group, which can be used to filter the IP address group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address_groups` -  The IP address groups.

  The [object](#address_groups_struct) structure is documented below.

<a name="address_groups_struct"></a>
The `address_groups` block supports:

* `addresses` -  IP address sets in an IP address group.
  Value range: a single IP address, IP address range, or CIDR block.

* `created_at` -  Time when the IP address group is created.

* `description` -  The supplementary information about the IP address group.

* `enterprise_project_id` -  Enterprise project ID.

* `id` -  IP address group ID, which uniquely identifies the IP address group.

* `ip_extra_set` -  IP addresses and their remarks in an IP address group.

  The [object](#address_groups_ip_extra_set_struct) structure is documented below.

* `ip_version` -  Whether it is an IPv4 or IPv6 address group.

* `max_capacity` -  Maximum number of entries in an address group,
  which limits the number of addresses that can be contained in an address group.

* `name` -  IP address group name.

* `status` -  The status of IP address group. Valid values are:
    + `NORMAL`: normal status.
    + `UPDATING`: updating.
    + `UPDATE_FAILED`: update failed.
  When the status of IP address group is `UPDATING`, the IP address group cannot be updated again.

* `status_message` -  The status details of IP address group.

* `updated_at` -  Time when the IP address group was last updated.

<a name="address_groups_ip_extra_set_struct"></a>
The `ip_extra_set` block supports:

* `ip` -  An IP address, IP address range, or CIDR block.

* `remarks` -  Provides supplementary information about the IP address, IP address range, or CIDR block.
