---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_acls_by_tags"
description: |-
  Use this data source to get a list of network ACLs by tags.
---

# huaweicloud_vpc_network_acls_by_tags

Use this data source to get a list of network ACLs by tags.

## Example Usage

```hcl
data "huaweicloud_vpc_network_acls_by_tags" "test" {
  tags {
    key    = "foo"
    values = ["bar"]
  }

  tags {
    key    = "key"
    values = ["value_1", "value_2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tags` - (Optional, List) Specifies the tags to filter to resources.
  The [tags](#tags) structure is documented below.

* `matches` - (Optional, List) Specifies the matches to filter to resources.
  The [matches](#matches) structure is documented below.

<a name="tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the values of the tag.

<a name="matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of the match. The value can be: **resource_name**.

* `value` - (Required, String) Specifies the value of the match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of the network ACLs found. The [resources](#resources) structure is documented below.

* `total_count` - The total count of the network ACLs found.

<a name="resources"></a>
The `resources` block supports:

* `resource_name` - The name of the network ACL.

* `resource_id` - The ID of the network ACL.

* `resource_detail` - The detail of the network ACL.

* `tags` - The tags which associated with the network ACL.
