---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_ips_by_tags"
description: |-
  Use this data source to get the list of transit IPs by tags within HuaweiCloud.
---

# huaweicloud_nat_private_transit_ips_by_tags

Use this data source to get the list of transit IPs by tags within HuaweiCloud.

## Example Usage

```hcl
variable "action" {}

data "huaweicloud_nat_private_transit_ips_by_tags" "test" {
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the action name.
  The valid values are **filter** and **count**.
  If `action` set to **filter**, indicates query the transit IP based on tags filtering conditions.
  If `action` set to **count**, indicates only query the total number of transit IPs.

* `tags` - (Optional, List) Specifies the resources to be queried contain tags list.
  
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies the resources to be queried contain any tag list.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the resources to be queried do not contain tags list.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the resources to be queried do not contain any tag list.
  The [not_tags_any](#tags_struct) structure is documented below.

-> For arguments above, include `tags`, `tags_any`, `not_tags`, `not_tags_any` have limits as follows:
  <br/>1. Each resource to be queried contains a maximum of `20` keys. Each tag key can have a maximum of `20`
  tag values. Each tag key must be unique, and each tag value in a tag must be unique.
  <br/>2. The structure body cannot be missing, and the key cannot be left blank or set to an empty string.
  <br/>3. Keys are in AND relationship (`tags`,`not_tags), and values in a tag are in OR relationship.
  <br/>4. Keys are in OR relationship (`tags_any`,`not_tags_any`), and values in a tag are in OR relationship.

* `matches` - (Optional, List) Specifies the search field.
  The tag `key` is the field to be matched, the value only can be **resource_name**.
  The tag `value` indicates the matched value.

  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tags key.
  The `key` cannot be left blank or be an empty string.
  The `key` cannot contain spaces.

* `values` - (Required, List) Specifies the list of tags values.
  The value can be an empty array but cannot be left blank. If values is an empty array, any value can be queried.
  The values are in the OR relationship.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the tag key used to search for resources.

* `value` - (Required, String) Specifies the tag value used to search for resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of transit IPs.

* `resources` - The list of transit IPs.
  The [resources](#transit_ips_struct) structure is documented below.

<a name="transit_ips_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name (Also is the transit IP address).

* `resource_detail` - The detail of the matched resources.
  The value is a resource object used for extension and is left blank by default.

* `tags` - The tags list.
  The [tags](#transit_ips_tags_struct) structure is documented below.

<a name="transit_ips_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
