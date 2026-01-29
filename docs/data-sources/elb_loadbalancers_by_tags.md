---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancers_by_tags"
description: |-
  Use this data source to get the list of ELB load balancers filtered by tags within HuaweiCloud.
---


# huaweicloud_elb_loadbalancers_by_tags

Use this data source to get the list of ELB load balancers filtered by tags within HuaweiCloud.

## Example Usage

```hcl
variable "action" {}

data "huaweicloud_elb_loadbalancers_by_tags" "test" {
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the action name. Possible values are **count** and **filter**.
  + **count**: Querying count of data.
  + **filter**: Querying details of data.

* `tags` - (Optional, List) Specifies the list of the tags to be queried.
  The [tags](#tags_struct) structure is documented below.

  -> 1. The tags contains a maximum of `20` keys. Each tag key can have a maximum of `20` tag values.
  <br/>2. The key cannot be left blank or set to an empty string.
  <br/>3. Each tag key must be unique, and each tag value in a tag must be unique.

* `matches` - (Optional, List) Specifies the search field.
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the resource tag.
  It contains a maximum of `128` unicode characters.

* `values` - (Required, List) Specifies the list of values corresponding to the key.
  
  -> 1. The tag value contains up to `255` unicode characters.
  <br/>2. The values in this list are in an OR relationship.
  <br/>3. If this list is left blank, it indicates that all values are included.
  <br/>4. If the value starts with (*), it indicates that fuzzy match is performed based on the value following (*).

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the matchs key.
  Currently, only **resource_name** for key is supported.

* `value` - (Required, String) Specifies the value corresponding to the key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of loadbalancers.

* `resources` - The list of loadbalancers.
  The [resources](#elb_loadbalancers_resources) structure is documented below.

<a name="elb_loadbalancers_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `super_resource_id` - The parent resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the resource.
  The value is a resource object used for extension. This parameter is left blank by default.

* `tags` - The tag list.
  The [tags](#elb_loadbalancers_tags) structure is documented below.

<a name="elb_loadbalancers_tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
