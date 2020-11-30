---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_secgroup

Use this data source to get the ID of an available HuaweiCloud security group.
This is an alternative to `huaweicloud_networking_secgroup_v2`

## Example Usage

```hcl
data "huaweicloud_networking_secgroup" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve security groups ids. If omitted, the
  `region` argument of the provider is used.

* `secgroup_id` - (Optional, String) The ID of the security group.

* `name` - (Optional, String) The name of the security group.

* `tenant_id` - (Optional, String) The owner of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `description`- The description of the security group.
